package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/crypto"
)

func TestRSAEncryptAndDecrypt(t *testing.T) {

	pub_key_string := `-----BEGIN RSA PUBLIC KEY-----
MIGJAoGBAMCYQQyPX45an+xB0Nrg5jA/DJlCkn2LJEHUYYPQIofVjTHtuPU0gWlq
1j13hXdVELUKH+VUEUu2dcwYDSZvSnnszAzFQs5xa3nxiI5Fqk5MFZgpwwKeSC1y
woL8WwdYzDL9TYQL2AUtMrRt3/Sw6FftXG/GslKgA5vz/eeVCbIhAgMBAAE=
-----END RSA PUBLIC KEY-----`
	priv_key_string := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDAmEEMj1+OWp/sQdDa4OYwPwyZQpJ9iyRB1GGD0CKH1Y0x7bj1
NIFpatY9d4V3VRC1Ch/lVBFLtnXMGA0mb0p57MwMxULOcWt58YiORapOTBWYKcMC
nkgtcsKC/FsHWMwy/U2EC9gFLTK0bd/0sOhX7VxvxrJSoAOb8/3nlQmyIQIDAQAB
AoGAce5fpL935q1jp45zr9eVNHtHx64jHJNgOnUZzrEkjDhfU2buoFeUKrlhzXDU
CnjO7lnz7mTh6mkgnECqHs99PUYNPKOKsU5beOJwqPNEMsVmY2wIUSXKp7013eXs
5v5Xdf7PDHu2o5SRT/SB6pQ1J5S1BsqCE+NDNXfRPFN/M1UCQQDsCQ/jl1n/HSUJ
uiFsTEIDYHB7JgFywmeSp0yZJn50EdV8byuL/vIQZPS0qXFTjML8fNhUPIbIvCJ/
Ha9IrTB7AkEA0OKPUNvsaCHD1V3YK+7Xz5jXc0PtEABslL8VbPdRHwFEsX6LWjIT
0wUlfJ9pRkmx8XGhr2buB7aP5ONcsmx7EwJBAKBaK7Q3f4mEWEQ6cjhrujEnFGNl
V3iKP+juxWgKMcBS2VE3CUOLiRHANEqEDpxvNYxomGLp17uJrHnlRc6+8f8CQH8x
FZdk8uTNepOXmyPVQa/1H2vedqGBwJwqZn99cPXyLcPujCgVyiB6R8NExjO4eBPO
32cQw+wKbEAxeaZji+UCQHN730WZ9GQhDIbXzxzqPBR6jmWQZ8w7czRQbLt2KR+Z
jXv37e+WoXSBWcRfi+MJx9uVabKcycf/Q2cmGoyN1c4=
-----END RSA PRIVATE KEY-----`
	pub_key, err := crypto.RSAReadPublicKey(pub_key_string)
	if err != nil {
		t.Fatal()
	}
	priv_key, err := crypto.RSAReadPrivateKey(priv_key_string)
	if err != nil {
		t.Fatal()
	}

	plain_text := "aaaaaaaaaaaaaaa"
	encrypt_text := ""

	encrypt_text, err = crypto.RSAEncrypt(plain_text, pub_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(encrypt_text)

	plain_text, err = crypto.RSADecrypt(encrypt_text, priv_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(plain_text)

	encrypt_text, err = crypto.RSAEncryptToBase64String(plain_text, pub_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(encrypt_text)

	plain_text, err = crypto.RSADecryptFromBase64String(encrypt_text, priv_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(plain_text)

	encrypt_text, err = crypto.RSAEncryptToHexString(plain_text, pub_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(encrypt_text)

	plain_text, err = crypto.RSADecryptFromHexString(encrypt_text, priv_key)
	if err != nil {
		t.Fatal()
	}
	fmt.Println(plain_text)

}
