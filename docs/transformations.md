## Transformations

### Template

A Template Transformation is able to render a text-based template with secrets previously
pulled from vaults. To use a template, add the following section to a configuration:

```yaml
transformations:
  - type: template
    in:
      - inputVar1
      - inputVar2
    out: outputVar
    spec:
      template: |
        other={{ .inputVar1 }}
        value={{ .inputVar2 }}
```

The above snippet renders the following output and stores it into `outputVar` for further processing,
given that `inputVar1` equals `some` and `inputVar2` equals `secret`.

```
other=some
value=secret
```

To add a Content Type, add to the spec part:

```yaml
    spec:
      contentType: text/plain
      template: |
        other={{ .inputVar1 }}
        value={{ .inputVar2 }}
```
