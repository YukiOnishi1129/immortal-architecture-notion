
# ModelsTemplateResponse

テンプレートレスポンス

## Properties

Name | Type
------------ | -------------
`id` | string
`name` | string
`ownerId` | string
`owner` | [ModelsAccountSummary](ModelsAccountSummary.md)
`fields` | [Array&lt;ModelsField&gt;](ModelsField.md)
`updatedAt` | Date
`isUsed` | boolean

## Example

```typescript
import type { ModelsTemplateResponse } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "name": null,
  "ownerId": null,
  "owner": null,
  "fields": null,
  "updatedAt": null,
  "isUsed": null,
} satisfies ModelsTemplateResponse

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsTemplateResponse
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


