
# ModelsNoteFilters

ノートフィルター（クエリパラメータ）

## Properties

Name | Type
------------ | -------------
`q` | string
`status` | [ModelsNoteStatus](ModelsNoteStatus.md)
`templateId` | string
`ownerId` | string

## Example

```typescript
import type { ModelsNoteFilters } from ''

// TODO: Update the object below with actual values
const example = {
  "q": null,
  "status": null,
  "templateId": null,
  "ownerId": null,
} satisfies ModelsNoteFilters

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsNoteFilters
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


