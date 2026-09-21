# NotesApi

All URIs are relative to *https://api.mini-notion.com*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**notesCreateNote**](NotesApi.md#notescreatenote) | **POST** /api/notes | Create note |
| [**notesDeleteNote**](NotesApi.md#notesdeletenote) | **DELETE** /api/notes/{noteId} | Delete note |
| [**notesGetNoteById**](NotesApi.md#notesgetnotebyid) | **GET** /api/notes/{noteId} | Get note by ID |
| [**notesListNotes**](NotesApi.md#noteslistnotes) | **GET** /api/notes | Get notes list |
| [**notesPublishNote**](NotesApi.md#notespublishnote) | **POST** /api/notes/{noteId}/publish | Publish note |
| [**notesUnpublishNote**](NotesApi.md#notesunpublishnote) | **POST** /api/notes/{noteId}/unpublish | Unpublish note |
| [**notesUpdateNote**](NotesApi.md#notesupdatenote) | **PUT** /api/notes/{noteId} | Update note |



## notesCreateNote

> ModelsNoteResponse notesCreateNote(modelsCreateNoteRequest)

Create note

ノート作成

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesCreateNoteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // ModelsCreateNoteRequest
    modelsCreateNoteRequest: ...,
  } satisfies NotesCreateNoteRequest;

  try {
    const data = await api.notesCreateNote(body);
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
| **modelsCreateNoteRequest** | [ModelsCreateNoteRequest](ModelsCreateNoteRequest.md) |  | |

### Return type

[**ModelsNoteResponse**](ModelsNoteResponse.md)

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


## notesDeleteNote

> ModelsSuccessResponse notesDeleteNote(noteId, ownerId)

Delete note

ノート削除

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesDeleteNoteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string
    noteId: noteId_example,
    // string | 所有者ID（権限チェック用）
    ownerId: ownerId_example,
  } satisfies NotesDeleteNoteRequest;

  try {
    const data = await api.notesDeleteNote(body);
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
| **noteId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者ID（権限チェック用） | [Defaults to `undefined`] |

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


## notesGetNoteById

> ModelsNoteResponse notesGetNoteById(noteId)

Get note by ID

ノート詳細取得

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesGetNoteByIdRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string
    noteId: noteId_example,
  } satisfies NotesGetNoteByIdRequest;

  try {
    const data = await api.notesGetNoteById(body);
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
| **noteId** | `string` |  | [Defaults to `undefined`] |

### Return type

[**ModelsNoteResponse**](ModelsNoteResponse.md)

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


## notesListNotes

> Array&lt;ModelsNoteResponse&gt; notesListNotes(q, status, templateId, ownerId)

Get notes list

ノート一覧取得

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesListNotesRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string | タイトルキーワード検索 (optional)
    q: q_example,
    // ModelsNoteStatus | ステータスフィルター (optional)
    status: ...,
    // string | テンプレートIDフィルター (optional)
    templateId: templateId_example,
    // string | 所有者IDフィルター (optional)
    ownerId: ownerId_example,
  } satisfies NotesListNotesRequest;

  try {
    const data = await api.notesListNotes(body);
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
| **q** | `string` | タイトルキーワード検索 | [Optional] [Defaults to `undefined`] |
| **status** | `ModelsNoteStatus` | ステータスフィルター | [Optional] [Defaults to `undefined`] [Enum: Draft, Publish] |
| **templateId** | `string` | テンプレートIDフィルター | [Optional] [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者IDフィルター | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsNoteResponse&gt;**](ModelsNoteResponse.md)

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


## notesPublishNote

> ModelsNoteResponse notesPublishNote(noteId, ownerId)

Publish note

ノート公開

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesPublishNoteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string
    noteId: noteId_example,
    // string | 所有者ID（公開権限チェック用）
    ownerId: ownerId_example,
  } satisfies NotesPublishNoteRequest;

  try {
    const data = await api.notesPublishNote(body);
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
| **noteId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者ID（公開権限チェック用） | [Defaults to `undefined`] |

### Return type

[**ModelsNoteResponse**](ModelsNoteResponse.md)

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


## notesUnpublishNote

> ModelsNoteResponse notesUnpublishNote(noteId, ownerId)

Unpublish note

ノート公開取り消し

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesUnpublishNoteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string
    noteId: noteId_example,
    // string | 所有者ID（公開権限チェック用）
    ownerId: ownerId_example,
  } satisfies NotesUnpublishNoteRequest;

  try {
    const data = await api.notesUnpublishNote(body);
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
| **noteId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者ID（公開権限チェック用） | [Defaults to `undefined`] |

### Return type

[**ModelsNoteResponse**](ModelsNoteResponse.md)

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


## notesUpdateNote

> ModelsNoteResponse notesUpdateNote(noteId, ownerId, modelsUpdateNoteRequest)

Update note

ノート更新

### Example

```ts
import {
  Configuration,
  NotesApi,
} from '';
import type { NotesUpdateNoteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new NotesApi();

  const body = {
    // string
    noteId: noteId_example,
    // string | 所有者ID（権限チェック用）
    ownerId: ownerId_example,
    // ModelsUpdateNoteRequest
    modelsUpdateNoteRequest: ...,
  } satisfies NotesUpdateNoteRequest;

  try {
    const data = await api.notesUpdateNote(body);
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
| **noteId** | `string` |  | [Defaults to `undefined`] |
| **ownerId** | `string` | 所有者ID（権限チェック用） | [Defaults to `undefined`] |
| **modelsUpdateNoteRequest** | [ModelsUpdateNoteRequest](ModelsUpdateNoteRequest.md) |  | |

### Return type

[**ModelsNoteResponse**](ModelsNoteResponse.md)

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

