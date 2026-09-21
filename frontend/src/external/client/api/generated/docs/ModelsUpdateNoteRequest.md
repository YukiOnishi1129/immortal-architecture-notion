
# ModelsUpdateNoteRequest

ノート更新リクエスト

## Properties

Name | Type
------------ | -------------
`id` | string
`title` | string
`sections` | [Array&lt;ModelsUpdateSectionRequest&gt;](ModelsUpdateSectionRequest.md)

## Example

```typescript
import type { ModelsUpdateNoteRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "title": null,
  "sections": null,
} satisfies ModelsUpdateNoteRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsUpdateNoteRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


