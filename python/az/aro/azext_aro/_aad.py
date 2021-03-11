# Copyright (c) Microsoft Corporation.
# Licensed under the Apache License 2.0.

import datetime
import time
import uuid

from azure.cli.core._profile import Profile
from azure.cli.core.commands.client_factory import configure_common_settings
from azure.graphrbac import GraphRbacManagementClient
from azure.graphrbac.models import ApplicationCreateParameters
from azure.graphrbac.models import GraphErrorException
from azure.graphrbac.models import PasswordCredential
from azure.graphrbac.models import ServicePrincipalCreateParameters
from knack.log import get_logger

logger = get_logger(__name__)


class AADManager:
    MANAGED_APP_PREFIX = 'https://az.aro.azure.com/'

    def __init__(self, cli_ctx):
        profile = Profile(cli_ctx=cli_ctx)
        credentials, _, tenant_id = profile.get_login_credentials(
            resource=cli_ctx.cloud.endpoints.active_directory_graph_resource_id)
        self.client = GraphRbacManagementClient(
            credentials, tenant_id, base_url=cli_ctx.cloud.endpoints.active_directory_graph_resource_id)
        configure_common_settings(cli_ctx, self.client)

    def create_application(self, display_name):
        password = uuid.uuid4()
        end_date = datetime.datetime(2299, 12, 31, tzinfo=datetime.timezone.utc)

        app = self.client.applications.create(ApplicationCreateParameters(
            display_name=display_name,
            identifier_uris=[
                self.MANAGED_APP_PREFIX + str(uuid.uuid4()),
            ],
            password_credentials=[
                PasswordCredential(
                    end_date=end_date,
                    value=password,
                ),
            ],
        ))

        return app, password

    def get_service_principal(self, app_id):
        sps = list(self.client.service_principals.list(
            filter="appId eq '%s'" % app_id))
        if sps:
            return sps[0]
        return None

    def get_application_by_client_id(self, client_id):
        apps = list(self.client.applications.list(
            filter="appId eq '%s'" % client_id))
        if apps:
            return apps[0]
        return None

    def create_service_principal(self, app_id):
        max_retries = 3
        retries = 0
        while True:
            try:
                return self.client.service_principals.create(
                    ServicePrincipalCreateParameters(app_id=app_id))
            except GraphErrorException as ex:
                if retries >= max_retries:
                    raise
                retries += 1
                logger.warning("%s; retry %d of %d", ex, retries, max_retries)
                time.sleep(10)

    def refresh_application_credentials(self, object_id):
        password = uuid.uuid4()
        key_id = uuid.uuid4()
        start_date = datetime.datetime.utcnow()
        end_date = datetime.datetime(2299, 12, 31, tzinfo=datetime.timezone.utc)

        try:
            credentials = list(self.client.applications.list_password_credentials(object_id))
        except GraphErrorException as e:
            raise logger.error(e.message)

        # when appending credentials ALL fields must be present, otherwise
        # azure gives ambiguous errors about not being able to update old keys
        credentials.append(PasswordCredential(
            custom_key_identifier=str(start_date).encode(),
            key_id=str(key_id),
            start_date=start_date,
            end_date=end_date,
            value=password))

        self.client.applications.update_password_credentials(object_id, credentials)

        return password
