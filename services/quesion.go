package services

import "github.com/miekg/dns"

const (
	notIPQuery = 0
	_IP4Query  = 4
	_IP6Query  = 6
)

type Question struct {
	qname  string
	qtype  string
	qclass string
	q      dns.Question
}

func NewQuestion(q dns.Question) *Question {
	return &Question{
		q:      q,
		qname:  unFqdn(q.Name),
		qtype:  dns.TypeToString[q.Qtype],
		qclass: dns.ClassToString[q.Qclass]}
}

func (q *Question) String() string {
	return q.qname + " " + q.qclass + " " + q.qtype
}
func (q *Question) QueryType() int {
	if q.q.Qclass != dns.ClassINET {
		return notIPQuery
	}
	switch q.q.Qtype {
	case dns.TypeA:
		return _IP4Query
	case dns.TypeAAAA:
		return _IP6Query
	default:
		return notIPQuery
	}
}
func unFqdn(s string) string {
	if dns.IsFqdn(s) {
		return s[:len(s)-1]
	}
	return s
}
