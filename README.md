## TODO

### Routes

- [x] Create GET /contacts/ route
- [x] Create POST /contacts/route
- [ ] Create PUT /contacts/ route
- [ ] Create PATCH /contacts/ route
- [ ] Create DELETE /contacts/ route

### Fields

- [x] Create foreign tables to contain contacts' emails and phone numbers

### Improvements

- [ ] Request validation
- [ ] Change from go sql driver to an ORM, like [GORM](https://github.com/jinzhu/gorm)
- [ ] Try to send logs to somewhere ([Datadog?](https://www.datadoghq.com/))
- [ ] Use redis to cache anything
