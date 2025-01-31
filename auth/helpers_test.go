package auth

import (
	"crypto/elliptic"
	"io/ioutil"
	"math/big"
	"regexp"
	"testing"

	"github.com/sethgrid/pester"
	"github.com/stretchr/testify/assert"
)

func TestRetrievingIAPJSONWebKeys(t *testing.T) {
	t.Run("ReturnsKeyByKeyID", func(t *testing.T) {

		// act (if fails get new kid from https://www.gstatic.com/iap/verify/public_key-jwk and update expectancies until it works)
		publicKey, err := GetCachedIAPJWK("f9R3yg")

		if assert.Nil(t, err) {
			assert.Equal(t, elliptic.P256(), publicKey.Curve)

			expectedX := new(big.Int)
			expectedX, _ = expectedX.SetString("33754992528993959342082873952071099444905807959681776349240807143574023195992", 10)

			if assert.Equal(t, expectedX, publicKey.X) {

				expectedY := new(big.Int)
				expectedY, _ = expectedY.SetString("30017756976983295626595109856839943719662421701617989535808220756803905010317", 10)

				assert.Equal(t, expectedY, publicKey.Y)
			}
		}
	})
}

func TestRetrievingGoogleJSONWebKeys(t *testing.T) {
	t.Run("ReturnsKeyByKeyID", func(t *testing.T) {

		// get kid from https://www.googleapis.com/oauth2/v3/certs
		response, err := pester.Get("https://www.googleapis.com/oauth2/v3/certs")
		if !assert.Nil(t, err, "Did not expect error %v", err) {
			return
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if !assert.Nil(t, err, "Did not expect error %v", err) {
			return
		}

		re := regexp.MustCompile(`"kid": "([a-z0-9]+)"`)
		match := re.FindStringSubmatch(string(body))

		if !assert.Equal(t, 2, len(match)) {
			return
		}

		kid := match[1]

		// act
		publicKey, err := GetCachedGoogleJWK(kid)

		if assert.Nil(t, err) {
			//expectedN, _ := new(big.Int).SetString("22883553494265264849962968666657907504805623544973021658037520392026414912908107935748030728083495882081762652678496550897864167563163931576062199151755602975712328787269312391349384940420444398673538634337720031141468560593835881691536106533135811050212245110259374057269172352043114432466466695441690457594821465498358453562523710809436428389580788674957453154488942814542128887082850773825912428357703060544920653986000150074585742341414869466318529599631736263936716538597041745264999951488037842637445133870248707789451035542980932611156149133030616737479007834220794751241211944945859836084724383454731144769171", 10)

			//if assert.Equal(t, expectedN, publicKey.N) {

			expectedY := 65537
			assert.Equal(t, expectedY, publicKey.E)
			//}
		}
	})
}
