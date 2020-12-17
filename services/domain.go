package services

type Domain struct {
	Domain string `json:"domain" valid:"dns,required"`
	IP     string `json:"ip" valid:"ip,required"`
	Value  string `json:"value"`
}
