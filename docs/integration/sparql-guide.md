# SPARQL Guide — Querying the Ontology

EduTrack exposes a read-only SPARQL endpoint that proxies to the VEDO Hub
ontology. It answers `SELECT`, `ASK`, `DESCRIBE` and `CONSTRUCT` queries.

## Endpoint

```
GET /api/v1/sparql?query=<url-encoded-sparql>
Authorization: Bearer <token>
```

The response is `application/sparql-results+json`:

```json
{
  "head": { "vars": ["s", "label"] },
  "results": {
    "bindings": [
      { "s": {"type": "uri", "value": "https://...#Class"}, "label": {"type": "literal", "value": "Class"} }
    ]
  },
  "truncated": false
}
```

## Common queries

List all classes:

```sparql
SELECT ?s ?label WHERE { ?s a <https://www.w3.org/2002/07/owl#Class> }
```

Check whether a class exists:

```sparql
ASK WHERE { <https://example.org/Module> a <https://www.w3.org/2002/07/owl#Class> }
```

Get a class's properties:

```sparql
SELECT ?p WHERE { ?s a <https://www.w3.org/2002/07/owl#Class> . ?s ?p ?o }
```

## Restrictions

- **Read-only:** mutating forms (`INSERT`, `DELETE`, `LOAD`, `CLEAR`,
  `CREATE`, `DROP`) are rejected with `403 Forbidden`.
- **Empty/malformed query:** `400 Bad Request`.
- **Result truncation:** results are clipped at 10 000 rows and flagged with
  `"truncated": true`.
- **Rate limit:** 10 req/min per user (burst 2). Exceeding it returns
  `429 Too Many Requests` with a `Retry-After` header.
- **Timeout:** queries longer than 30s return `504 Gateway Timeout`.

## curl

```bash
QUERY='SELECT ?s ?label WHERE { ?s a <https://www.w3.org/2002/07/owl#Class> }'
curl -G http://localhost:8080/api/v1/sparql \
  -H "Authorization: Bearer $TOKEN" \
  --data-urlencode "query=$QUERY" | jq
```

## See also

- `examples/curl-examples.sh` — runnable examples
- `quickstart.md` — get started in 5 minutes
