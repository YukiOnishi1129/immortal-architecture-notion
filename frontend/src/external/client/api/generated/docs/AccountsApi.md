# AccountsApi

All URIs are relative to *https://api.mini-notion.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**accountsCreateOrGetAccount**](AccountsApi.md#accountscreateorgetaccount) | **POST** /api/accounts/auth | Create or get account via OAuth |
| [**accountsGetAccountByEmail**](AccountsApi.md#accountsgetaccountbyemail) | **GET** /api/accounts/by-email | Get account by email |
| [**accountsGetAccountById**](AccountsApi.md#accountsgetaccountbyid) | **GET** /api/accounts/{accountId} | Get account by ID |
| [**accountsGetCurrentAccount**](AccountsApi.md#accountsgetcurrentaccount) | **GET** /api/accounts/me | Get current account |



## accountsCreateOrGetAccount

> ModelsAccountResponse accountsCreateOrGetAccount(modelsCreateOrGetAccountRequest)

Create or get account via OAuth

OAuth認証（内部処理）

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '';
import type { AccountsCreateOrGetAccountRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AccountsApi();

  const body = {
    // ModelsCreateOrGetAccountRequest
    modelsCreateOrGetAccountRequest: ...,
  } satisfies AccountsCreateOrGetAccountRequest;

  try {
    const data = await api.accountsCreateOrGetAccount(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **modelsCreateOrGetAccountRequest** | [ModelsCreateOrGetAccountRequest](ModelsCreateOrGetAccountRequest.md) |  | |

### Return type

[**ModelsAccountResponse**](ModelsAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The request has succeeded. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## accountsGetAccountByEmail

> ModelsAccountResponse accountsGetAccountByEmail(email)

Get account by email

メールアドレスでアカウント取得

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '';
import type { AccountsGetAccountByEmailRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AccountsApi();

  const body = {
    // string
    email: email_example,
  } satisfies AccountsGetAccountByEmailRequest;

  try {
    const data = await api.accountsGetAccountByEmail(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **email** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ModelsAccountResponse**](ModelsAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The request has succeeded. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## accountsGetAccountById

> ModelsAccountResponse accountsGetAccountById(accountId)

Get account by ID

アカウント詳細取得

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '';
import type { AccountsGetAccountByIdRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AccountsApi();

  const body = {
    // string
    accountId: accountId_example,
  } satisfies AccountsGetAccountByIdRequest;

  try {
    const data = await api.accountsGetAccountById(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **accountId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ModelsAccountResponse**](ModelsAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The request has succeeded. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## accountsGetCurrentAccount

> ModelsAccountResponse accountsGetCurrentAccount()

Get current account

現在のアカウント取得

### Example

```ts
import {
  Configuration,
  AccountsApi,
} from '';
import type { AccountsGetCurrentAccountRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new AccountsApi();

  try {
    const data = await api.accountsGetCurrentAccount();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**ModelsAccountResponse**](ModelsAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | The request has succeeded. |  -  |
| **0** | An unexpected error response. |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

