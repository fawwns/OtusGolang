//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov",
	"Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com",
"Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com",
"Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net",
"Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com",
"Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStat_InvalidData(t *testing.T) {
	data := `{"Email":"TEST@Example.COM"}
{"Email":"broken_json"
{"Email":"user@Another.org"}
`

	result, err := GetDomainStat(bytes.NewBufferString(data), "com")
	require.NoError(t, err)

	require.Equal(t, DomainStat{"example.com": 1}, result)
}

func BenchmarkGetDomainStat(b *testing.B) {
	data := `
{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov",
"Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Brian Olson","Username":"non_quia_id","Email":"FrancesEllis@Quinu.edu",
"Phone":"237-75-34","Password":"cmEPhX8","Address":"Butterfield Junction 74"}
{"Id":3,"Name":"Justin Oliver Jr. Sr.","Username":"oPerez","Email":"MelissaGutierrez@Twinte.gov",
"Phone":"106-05-18","Password":"f00GKr9i","Address":"Oak Valley Lane 19"}
{"Id":4,"Name":"Joan Evans","Username":"eaque_nobis","Email":"ElmerBishop@Buzzshare.com",
"Phone":"762-53-33","Password":"gT4nTqk","Address":"North Pine Road 88"}
{"Id":5,"Name":"Paula Webb","Username":"et_quibusdam","Email":"CalebGarcia@Browsedrive.gov",
"Phone":"293-19-91","Password":"cK1LxYw","Address":"South Sunset Blvd 2"}
{"Id":6,"Name":"Billy Craig","Username":"dolor_voluptatem","Email":"MargieSpencer@Quinu.edu",
"Phone":"8-268-228-93-45","Password":"kEmHzJf","Address":"East Harbor Lane 41"}
{"Id":7,"Name":"Norma Adams","Username":"ullam_esse","Email":"JeanetteMurray@Twinte.gov",
"Phone":"598-61-60","Password":"yKzJrTx","Address":"Forest Hill 11"}
{"Id":8,"Name":"Bobby Mitchell","Username":"quod_ab","Email":"SaraClark@Browsedrive.gov",
"Phone":"4-987-231-77-34","Password":"KpXbUfq","Address":"Maplewood Court 99"}
{"Id":9,"Name":"Jean Knight","Username":"suscipit_fuga","Email":"AnthonyHarrison@Quinu.edu",
"Phone":"543-79-21","Password":"xWgTzqN","Address":"Old Mill Street 3"}
{"Id":10,"Name":"Roger Taylor","Username":"a_fugit","Email":"HelenMartinez@Twinte.gov",
"Phone":"239-91-53","Password":"UpQnLvr","Address":"Lakeshore Drive 67"}
`

	for i := 0; i < b.N; i++ {
		_, _ = GetDomainStat(bytes.NewBufferString(data), "com")
	}
}
