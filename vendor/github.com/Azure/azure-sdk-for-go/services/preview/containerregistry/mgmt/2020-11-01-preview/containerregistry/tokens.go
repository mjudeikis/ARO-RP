package containerregistry

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/validation"
	"github.com/Azure/go-autorest/tracing"
	"net/http"
)

// TokensClient is the client for the Tokens methods of the Containerregistry service.
type TokensClient struct {
	BaseClient
}

// NewTokensClient creates an instance of the TokensClient client.
func NewTokensClient(subscriptionID string) TokensClient {
	return NewTokensClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewTokensClientWithBaseURI creates an instance of the TokensClient client using a custom endpoint.  Use this when
// interacting with an Azure cloud that uses a non-standard base URI (sovereign clouds, Azure stack).
func NewTokensClientWithBaseURI(baseURI string, subscriptionID string) TokensClient {
	return TokensClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Create creates a token for a container registry with the specified parameters.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// registryName - the name of the container registry.
// tokenName - the name of the token.
// tokenCreateParameters - the parameters for creating a token.
func (client TokensClient) Create(ctx context.Context, resourceGroupName string, registryName string, tokenName string, tokenCreateParameters Token) (result TokensCreateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.Create")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: registryName,
			Constraints: []validation.Constraint{{Target: "registryName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "registryName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "registryName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]*$`, Chain: nil}}},
		{TargetValue: tokenName,
			Constraints: []validation.Constraint{{Target: "tokenName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "tokenName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "tokenName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("containerregistry.TokensClient", "Create", err.Error())
	}

	req, err := client.CreatePreparer(ctx, resourceGroupName, registryName, tokenName, tokenCreateParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Create", nil, "Failure preparing request")
		return
	}

	result, err = client.CreateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Create", nil, "Failure sending request")
		return
	}

	return
}

// CreatePreparer prepares the Create request.
func (client TokensClient) CreatePreparer(ctx context.Context, resourceGroupName string, registryName string, tokenName string, tokenCreateParameters Token) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"registryName":      autorest.Encode("path", registryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tokenName":         autorest.Encode("path", tokenName),
	}

	const APIVersion = "2020-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/tokens/{tokenName}", pathParameters),
		autorest.WithJSON(tokenCreateParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// CreateSender sends the Create request. The method will close the
// http.Response Body if it receives an error.
func (client TokensClient) CreateSender(req *http.Request) (future TokensCreateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = func(client TokensClient) (t Token, err error) {
		var done bool
		done, err = future.DoneWithContext(context.Background(), client)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerregistry.TokensCreateFuture", "Result", future.Response(), "Polling failure")
			return
		}
		if !done {
			err = azure.NewAsyncOpIncompleteError("containerregistry.TokensCreateFuture")
			return
		}
		sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
		t.Response.Response, err = future.GetResult(sender)
		if t.Response.Response == nil && err == nil {
			err = autorest.NewErrorWithError(err, "containerregistry.TokensCreateFuture", "Result", nil, "received nil response and error")
		}
		if err == nil && t.Response.Response.StatusCode != http.StatusNoContent {
			t, err = client.CreateResponder(t.Response.Response)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerregistry.TokensCreateFuture", "Result", t.Response.Response, "Failure responding to request")
			}
		}
		return
	}
	return
}

// CreateResponder handles the response to the Create request. The method always
// closes the http.Response Body.
func (client TokensClient) CreateResponder(resp *http.Response) (result Token, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Delete deletes a token from a container registry.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// registryName - the name of the container registry.
// tokenName - the name of the token.
func (client TokensClient) Delete(ctx context.Context, resourceGroupName string, registryName string, tokenName string) (result TokensDeleteFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.Delete")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: registryName,
			Constraints: []validation.Constraint{{Target: "registryName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "registryName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "registryName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]*$`, Chain: nil}}},
		{TargetValue: tokenName,
			Constraints: []validation.Constraint{{Target: "tokenName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "tokenName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "tokenName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("containerregistry.TokensClient", "Delete", err.Error())
	}

	req, err := client.DeletePreparer(ctx, resourceGroupName, registryName, tokenName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Delete", nil, "Failure preparing request")
		return
	}

	result, err = client.DeleteSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Delete", nil, "Failure sending request")
		return
	}

	return
}

// DeletePreparer prepares the Delete request.
func (client TokensClient) DeletePreparer(ctx context.Context, resourceGroupName string, registryName string, tokenName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"registryName":      autorest.Encode("path", registryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tokenName":         autorest.Encode("path", tokenName),
	}

	const APIVersion = "2020-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsDelete(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/tokens/{tokenName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// DeleteSender sends the Delete request. The method will close the
// http.Response Body if it receives an error.
func (client TokensClient) DeleteSender(req *http.Request) (future TokensDeleteFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = func(client TokensClient) (ar autorest.Response, err error) {
		var done bool
		done, err = future.DoneWithContext(context.Background(), client)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerregistry.TokensDeleteFuture", "Result", future.Response(), "Polling failure")
			return
		}
		if !done {
			err = azure.NewAsyncOpIncompleteError("containerregistry.TokensDeleteFuture")
			return
		}
		ar.Response = future.Response()
		return
	}
	return
}

// DeleteResponder handles the response to the Delete request. The method always
// closes the http.Response Body.
func (client TokensClient) DeleteResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent),
		autorest.ByClosing())
	result.Response = resp
	return
}

// Get gets the properties of the specified token.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// registryName - the name of the container registry.
// tokenName - the name of the token.
func (client TokensClient) Get(ctx context.Context, resourceGroupName string, registryName string, tokenName string) (result Token, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.Get")
		defer func() {
			sc := -1
			if result.Response.Response != nil {
				sc = result.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: registryName,
			Constraints: []validation.Constraint{{Target: "registryName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "registryName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "registryName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]*$`, Chain: nil}}},
		{TargetValue: tokenName,
			Constraints: []validation.Constraint{{Target: "tokenName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "tokenName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "tokenName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("containerregistry.TokensClient", "Get", err.Error())
	}

	req, err := client.GetPreparer(ctx, resourceGroupName, registryName, tokenName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Get", resp, "Failure responding to request")
		return
	}

	return
}

// GetPreparer prepares the Get request.
func (client TokensClient) GetPreparer(ctx context.Context, resourceGroupName string, registryName string, tokenName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"registryName":      autorest.Encode("path", registryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tokenName":         autorest.Encode("path", tokenName),
	}

	const APIVersion = "2020-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/tokens/{tokenName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client TokensClient) GetSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client TokensClient) GetResponder(resp *http.Response) (result Token, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List lists all the tokens for the specified container registry.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// registryName - the name of the container registry.
func (client TokensClient) List(ctx context.Context, resourceGroupName string, registryName string) (result TokenListResultPage, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.List")
		defer func() {
			sc := -1
			if result.tlr.Response.Response != nil {
				sc = result.tlr.Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: registryName,
			Constraints: []validation.Constraint{{Target: "registryName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "registryName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "registryName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("containerregistry.TokensClient", "List", err.Error())
	}

	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, registryName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.tlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "List", resp, "Failure sending request")
		return
	}

	result.tlr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "List", resp, "Failure responding to request")
		return
	}
	if result.tlr.hasNextLink() && result.tlr.IsEmpty() {
		err = result.NextWithContext(ctx)
		return
	}

	return
}

// ListPreparer prepares the List request.
func (client TokensClient) ListPreparer(ctx context.Context, resourceGroupName string, registryName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"registryName":      autorest.Encode("path", registryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
	}

	const APIVersion = "2020-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/tokens", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client TokensClient) ListSender(req *http.Request) (*http.Response, error) {
	return client.Send(req, azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client TokensClient) ListResponder(resp *http.Response) (result TokenListResult, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client TokensClient) listNextResults(ctx context.Context, lastResults TokenListResult) (result TokenListResult, err error) {
	req, err := lastResults.tokenListResultPreparer(ctx)
	if err != nil {
		return result, autorest.NewErrorWithError(err, "containerregistry.TokensClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "containerregistry.TokensClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client TokensClient) ListComplete(ctx context.Context, resourceGroupName string, registryName string) (result TokenListResultIterator, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.List")
		defer func() {
			sc := -1
			if result.Response().Response.Response != nil {
				sc = result.page.Response().Response.Response.StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	result.page, err = client.List(ctx, resourceGroupName, registryName)
	return
}

// Update updates a token with the specified parameters.
// Parameters:
// resourceGroupName - the name of the resource group to which the container registry belongs.
// registryName - the name of the container registry.
// tokenName - the name of the token.
// tokenUpdateParameters - the parameters for updating a token.
func (client TokensClient) Update(ctx context.Context, resourceGroupName string, registryName string, tokenName string, tokenUpdateParameters TokenUpdateParameters) (result TokensUpdateFuture, err error) {
	if tracing.IsEnabled() {
		ctx = tracing.StartSpan(ctx, fqdn+"/TokensClient.Update")
		defer func() {
			sc := -1
			if result.FutureAPI != nil && result.FutureAPI.Response() != nil {
				sc = result.FutureAPI.Response().StatusCode
			}
			tracing.EndSpan(ctx, sc, err)
		}()
	}
	if err := validation.Validate([]validation.Validation{
		{TargetValue: resourceGroupName,
			Constraints: []validation.Constraint{{Target: "resourceGroupName", Name: validation.MinLength, Rule: 1, Chain: nil}}},
		{TargetValue: registryName,
			Constraints: []validation.Constraint{{Target: "registryName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "registryName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "registryName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9]*$`, Chain: nil}}},
		{TargetValue: tokenName,
			Constraints: []validation.Constraint{{Target: "tokenName", Name: validation.MaxLength, Rule: 50, Chain: nil},
				{Target: "tokenName", Name: validation.MinLength, Rule: 5, Chain: nil},
				{Target: "tokenName", Name: validation.Pattern, Rule: `^[a-zA-Z0-9-]*$`, Chain: nil}}}}); err != nil {
		return result, validation.NewError("containerregistry.TokensClient", "Update", err.Error())
	}

	req, err := client.UpdatePreparer(ctx, resourceGroupName, registryName, tokenName, tokenUpdateParameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Update", nil, "Failure preparing request")
		return
	}

	result, err = client.UpdateSender(req)
	if err != nil {
		err = autorest.NewErrorWithError(err, "containerregistry.TokensClient", "Update", nil, "Failure sending request")
		return
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client TokensClient) UpdatePreparer(ctx context.Context, resourceGroupName string, registryName string, tokenName string, tokenUpdateParameters TokenUpdateParameters) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"registryName":      autorest.Encode("path", registryName),
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"tokenName":         autorest.Encode("path", tokenName),
	}

	const APIVersion = "2020-11-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerRegistry/registries/{registryName}/tokens/{tokenName}", pathParameters),
		autorest.WithJSON(tokenUpdateParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client TokensClient) UpdateSender(req *http.Request) (future TokensUpdateFuture, err error) {
	var resp *http.Response
	resp, err = client.Send(req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		return
	}
	var azf azure.Future
	azf, err = azure.NewFutureFromResponse(resp)
	future.FutureAPI = &azf
	future.Result = func(client TokensClient) (t Token, err error) {
		var done bool
		done, err = future.DoneWithContext(context.Background(), client)
		if err != nil {
			err = autorest.NewErrorWithError(err, "containerregistry.TokensUpdateFuture", "Result", future.Response(), "Polling failure")
			return
		}
		if !done {
			err = azure.NewAsyncOpIncompleteError("containerregistry.TokensUpdateFuture")
			return
		}
		sender := autorest.DecorateSender(client, autorest.DoRetryForStatusCodes(client.RetryAttempts, client.RetryDuration, autorest.StatusCodesForRetry...))
		t.Response.Response, err = future.GetResult(sender)
		if t.Response.Response == nil && err == nil {
			err = autorest.NewErrorWithError(err, "containerregistry.TokensUpdateFuture", "Result", nil, "received nil response and error")
		}
		if err == nil && t.Response.Response.StatusCode != http.StatusNoContent {
			t, err = client.UpdateResponder(t.Response.Response)
			if err != nil {
				err = autorest.NewErrorWithError(err, "containerregistry.TokensUpdateFuture", "Result", t.Response.Response, "Failure responding to request")
			}
		}
		return
	}
	return
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client TokensClient) UpdateResponder(resp *http.Response) (result Token, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK, http.StatusCreated),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
