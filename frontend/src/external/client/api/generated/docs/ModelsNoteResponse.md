
# ModelsNoteResponse

ノートレスポンス

## Properties

Name | Type
------------ | -------------
`id` | string
`title` | string
`templateId` | string
`templateName` | string
`ownerId` | string
`owner` | [ModelsAccountSummary](ModelsAccountSummary.md)
`status` | [ModelsNoteStatus](ModelsNoteStatus.md)
`sections` | [Array&lt;ModelsSection&gt;](ModelsSection.md)
`createdAt` | Date
`updatedAt` | Date

## Example

```typescript
import type { ModelsNoteResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "title": null,
  "templateId": null,
  "templateName": null,
  "ownerId": null,
  "owner": null,
  "status": null,
  "sections": null,
  "createdAt": null,
  "updatedAt": null,
} satisfies ModelsNoteResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsNoteResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


