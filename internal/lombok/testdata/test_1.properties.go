package testdata

// properties for T
func (t *T) P1() string {
	return t.p1
}
func (t *T) GetP2() string {
	return t.p2
}
func (t *T) GetP3() string {
	return t.p3
}
func (t *T) SetP3(v string) {
	t.p3 = v
}
