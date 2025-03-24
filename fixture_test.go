package fixture_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gost-dom/fixture"
)

func TestOpenStreetMap(t *testing.T) {
	// This simple test demonstrates testing an place search component that uses
	// openstreetmap.
	//
	// It uses the Server from httptest to create an ephemeral server with
	// controlled endpoints that permit the test to simulate OSM responses

	// - HTTPHandlerFixture controls the endpoints using http.ServeMux
	// - HTTPServerFixture launches an ephemeral test server, serving the HTTPHandlerFixture's http handler
	// - OpenStreetMapFixture depends on on the HTTPServerFixture. It initializes OSM with the base URL from the server
	//
	// The details that the test cares about is that a specific mocked response
	// is converted into a specific structure.

	// The underlying philosophy is that the test should make explicit the
	// details that affect the outcome of the test, and hide the details that
	// are irrelevant for the outcome of the test.
	//
	// This test is about how a specific response is converted to a Go type, and
	// a recorded historic response is used as input. The test controls an http
	// ServeMux instance to setup the response, but the details of how the OSM
	// component is connected to that handler is hidden in the fixtures.

	fix, ctrl := fixture.Init(t, &struct {
		*OpenStreetMapFixture
		*HTTPHandlerFixture
		*HTTPServerFixture
	}{})
	ctrl.Setup()

	// Currently, cleanup is not supported, so we need to do this explicitly.
	t.Cleanup(fix.HTTPServerFixture.Close)
	fix.ServeMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})

	actual, err := fix.Search("dueodde")
	if err != nil {
		t.Errorf("Search returned an error: %v", err)
	}
	expected := []Result{
		{
			PlaceID:     135627303,
			Name:        "Dueodde",
			DisplayName: "Dueodde, Bornholm Regional Municipality, Capital Region of Denmark, Denmark",
		},
		{
			PlaceID:     136558536,
			Name:        "Dueodde",
			DisplayName: "Dueodde, Sirenevej, Bornholm Regional Municipality, Capital Region of Denmark, Denmark",
		},
		{
			PlaceID:     136968598,
			Name:        "Dueodde",
			DisplayName: "Dueodde, Udegårdsvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark",
		},
		{
			PlaceID:     134967504,
			Name:        "Dueodde Fyr",
			DisplayName: "Dueodde Fyr, Fyrvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark",
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Unexpected result\n  Expected: %v\n\n Actual: %v", expected, actual)
	}
}

type HTTPHandlerFixture struct {
	*http.ServeMux
}

func (f *HTTPHandlerFixture) Mux() *http.ServeMux {
	f.Setup()
	return f.ServeMux
}

func (f *HTTPHandlerFixture) Setup() {
	if f.ServeMux == nil {
		f.ServeMux = http.NewServeMux()
	}
}

type HTTPServerFixture struct {
	Handler *HTTPHandlerFixture
	*httptest.Server
}

func (f *HTTPServerFixture) Setup() {
	if f.Server == nil {
		f.Server = httptest.NewServer(f.Handler.ServeMux)
	}
}

type OpenStreetMapFixture struct {
	fixture.Fixture
	OpenStreetMap
	ServerFixture *HTTPServerFixture
}

func (f *OpenStreetMapFixture) Setup() {
	f.OpenStreetMap.BaseURL = f.ServerFixture.URL
}

const response = `[{"place_id":135627303,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":2960711590,"lat":"54.9902757","lon":"15.0759924","class":"place","type":"locality","place_rank":25,"importance":0.06672011825741969,"addresstype":"locality","name":"Dueodde","display_name":"Dueodde, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9802757","55.0002757","15.0659924","15.0859924"]},{"place_id":136558536,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":792653018,"lat":"54.9935299","lon":"15.0781951","class":"highway","type":"bus_stop","place_rank":30,"importance":5.345159075298722e-05,"addresstype":"highway","name":"Dueodde","display_name":"Dueodde, Sirenevej, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9934799","54.9935799","15.0781451","15.0782451"]},{"place_id":136968598,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":10006025636,"lat":"54.9954556","lon":"15.0381217","class":"tourism","type":"information","place_rank":30,"importance":9.99999999995449e-06,"addresstype":"tourism","name":"Dueodde","display_name":"Dueodde, Udegårdsvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9954056","54.9955056","15.0380717","15.0381717"]},{"place_id":134967504,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"way","osm_id":528547308,"lat":"54.991797500000004","lon":"15.074294876767677","class":"man_made","type":"lighthouse","place_rank":30,"importance":0.3797855310791376,"addresstype":"man_made","name":"Dueodde Fyr","display_name":"Dueodde Fyr, Fyrvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9917777","54.9918173","15.0742643","15.0743255"]}]`
