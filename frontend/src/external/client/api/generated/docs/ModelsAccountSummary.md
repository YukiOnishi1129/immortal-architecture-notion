
# ModelsAccountSummary

簡易アカウント情報（他のレスポンスに埋め込まれる）

## Properties

Name | Type
------------ | -------------
`id` | string
`firstName` | string
`lastName` | string
`thumbnail` | string

## Example

```typescript
import type { ModelsAccountSummary } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "firstName": null,
  "lastName": null,
  "thumbnail": null,
} satisfies ModelsAccountSummary

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsAccountSummary
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


