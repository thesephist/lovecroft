# Lovecroft

Lovecroft is a minimal mailing list manager supporting multiple mailing lists. It backs newsletters behind [thesephist.com](https://thesephist.com) and [Atypical Press](https://atypicalpress.com).

## API

### POST `/subscribe/{listName}`

```ts
{
    givenName: string,
    familyName: string,
    email: string,
}
```

### GET `/unsubscribe/{listName}/{token}`

### GET `/directory`

### GET `/list/{listName}`

### GET `/list-csv/{listName}`