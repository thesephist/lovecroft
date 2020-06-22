# Lovecroft

Lovecroft is a minimal mailing list manager supporting multiple mailing lists. It backs newsletters behind [thesephist.com](https://thesephist.com) and [Atypical Press](https://atypicalpress.com).

## API

### POST `/subscribe`

```ts
{
    givenName: string,
    familyName: string,
    email: string,
}
```

### GET `/unsubscribe`

Query parameters

- `token: string` unsubscribe token

### GET `/directory`

### GET `/list/{listName}`

### GET `/list-csv/{listName}`