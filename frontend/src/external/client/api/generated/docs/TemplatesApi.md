# TemplatesApi

All URIs are relative to *https://api.mini-notion.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**templatesCreateTemplate**](TemplatesApi.md#templatescreatetemplate) | **POST** /api/templates | Create template |
| [**templatesDeleteTemplate**](TemplatesApi.md#templatesdeletetemplate) | **DELETE** /api/templates/{templateId} | Delete template |
| [**templatesGetTemplateById**](TemplatesApi.md#templatesgettemplatebyid) | **GET** /api/templates/{templateId} | Get template by ID |
| [**templatesListTemplates**](TemplatesApi.md#templateslisttemplates) | **GET** /api/templates | Get templates list |
| [**templatesUpdateTemplate**](TemplatesApi.md#templatesupdatetemplate) | **PUT** /api/templates/{templateId} | Update template |



## templatesCreateTemplate

> ModelsTemplateResponse templatesCreateTemplate(modelsCreateTemplateRequest)

Create template

テンプレート作成

### Example

```ts
import {
  Configuration,
  TemplatesApi,
} from '';
import type { TemplatesCreateTemplateRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TemplatesApi();

  const body = {
    // ModelsCreateTemplateRequest
    modelsCreateTemplateRequest: ...,
  } satisfies TemplatesCreateTemplateRequest;

  try {
    const data = await api.templatesCreateTemplate(body);
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
| **modelsCreateTemplateRequest** | [ModelsCreateTemplateRequest](ModelsCreateTemplateRequest.md) |  | |

### Return type

[**ModelsTemplateResponse**](ModelsTemplateResponse.md)

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


## templatesDeleteTemplate

> ModelsSuccessResponse templatesDeleteTemplate(templateId, ownerId)

Delete template

テンプレート削除

### Example

```ts
import {
  Configuration,
  TemplatesApi,
} from '';
import type { TemplatesDeleteTemplateRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TemplatesApi();

  const body = {
    // string
    templateId: templateId_example,
    // string
    ownerId: ownerId_example,
  } satisfies TemplatesDeleteTemplateRequest;

  try {
    const data = await api.templatesDeleteTemplate(body);
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
| **templateId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ModelsSuccessResponse**](ModelsSuccessResponse.md)

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


## templatesGetTemplateById

> ModelsTemplateResponse templatesGetTemplateById(templateId)

Get template by ID

テンプレート詳細取得

### Example

```ts
import {
  Configuration,
  TemplatesApi,
} from '';
import type { TemplatesGetTemplateByIdRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TemplatesApi();

  const body = {
    // string
    templateId: templateId_example,
  } satisfies TemplatesGetTemplateByIdRequest;

  try {
    const data = await api.templatesGetTemplateById(body);
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
| **templateId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ModelsTemplateResponse**](ModelsTemplateResponse.md)

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


## templatesListTemplates

> Array&lt;ModelsTemplateResponse&gt; templatesListTemplates(q, ownerId)

Get templates list

テンプレート一覧取得

### Example

```ts
import {
  Configuration,
  TemplatesApi,
} from '';
import type { TemplatesListTemplatesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TemplatesApi();

  const body = {
    // string | テンプレート名のキーワード検索 (optional)
    q: q_example,
    // string | 所有者IDフィルター (optional)
    ownerId: ownerId_example,
  } satisfies TemplatesListTemplatesRequest;

  try {
    const data = await api.templatesListTemplates(body);
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
| **q** | `string` | テンプレート名のキーワード検索 | [Optional] [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者IDフィルター | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsTemplateResponse&gt;**](ModelsTemplateResponse.md)

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


## templatesUpdateTemplate

> ModelsTemplateResponse templatesUpdateTemplate(templateId, ownerId, modelsUpdateTemplateRequest)

Update template

テンプレート更新

### Example

```ts
import {
  Configuration,
  TemplatesApi,
} from '';
import type { TemplatesUpdateTemplateRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new TemplatesApi();

  const body = {
    // string
    templateId: templateId_example,
    // string
    ownerId: ownerId_example,
    // ModelsUpdateTemplateRequest
    modelsUpdateTemplateRequest: ...,
  } satisfies TemplatesUpdateTemplateRequest;

  try {
    const data = await api.templatesUpdateTemplate(body);
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
| **templateId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` |  | [Defaults to `undefined`] |
| **modelsUpdateTemplateRequest** | [ModelsUpdateTemplateRequest](ModelsUpdateTemplateRequest.md) |  | |

### Return type

[**ModelsTemplateResponse**](ModelsTemplateResponse.md)

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

