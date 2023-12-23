package query

// We want:
// Filter
// Sort
// Pagination
// Relation

// FILTER:
// , = and in the same field
// | = or in the same field
// ; = or in different fields
// : = and in different fields
// () = grouping
// [] = operator
// String value has to be between quotes
// url.com?filter=name[eq]value,otherfield[eq]othervalue
// url.com?filter=name[eq]value;

// SORT:
// url.com?sort=name[asc],otherfield[desc]
// url.com?sort=name:asc,otherfield:desc

// PAGINATION:
// url.com?page=1&limit=10

// RELATION:
// url.com?expand=relation,otherrelation
