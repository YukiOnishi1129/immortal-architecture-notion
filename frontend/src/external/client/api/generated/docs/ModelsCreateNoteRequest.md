
# ModelsCreateNoteRequest

ノート作成リクエスト

## Properties

Name | Type
------------ | -------------
`title` | string
`templateId` | string
`ownerId` | string
`sections` | [Array&lt;ModelsCreateSectionRequest&gt;](ModelsCreateSectionRequest.md)

## Example

```typescript
import type { ModelsCreateNoteRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "title": null,
  "templateId": null,
  "ownerId": null,
  "sections": null,
} satisfies ModelsCreateNoteRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsCreateNoteRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


