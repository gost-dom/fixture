package fixture_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gost-dom/fixture"
)

// This simple test demonstrates testing an place search component that uses
// openstreetmap.
//
// The underlying philosophy is that the test should make explicit the details
// that affect the outcome of the test, and hide the details that are irrelevant
// for the outcome of the test.
//
// The details that the test cares about is that a specific mocked response is
// converted into a specific Go structure. To be explicit about this, the test
// sets up canned responses for specific endpoints.
//
// The test uses httptest.Server to create an ephemeral server exposing the
// endpoints controlled by tests. This is a detail not relevant for the specific
// test and has been made implicit.
//
// The following three fixtures are used in the test:
//
// - HTTPHandlerFixture - create an new http.ServeMux
// - HTTPServerFixture - launches an ephemeral test server, serving the HTTPHandlerFixture's http handler
// - OpenStreetMapFixture - Initializes OSM, setting the BaseURL with a value exposed from HTTPServerFixture
//
// The test creates an inline struct to get the two components it needs, the
// OpenStreetMapFixture to exercise the code, and the HTTPHandlerFixture to
// setup a canned response.
//
// HTTPServerFixture that controls the server is a detail that is not visible in
// the test. Both setup and cleanup of the ephemeral HTTP server happens
// automatically
func TestOpenStreetMap(t *testing.T) {
	t.Parallel()

	fix, ctrl := fixture.Init(t, &struct {
		*OpenStreetMapFixture
		*HTTPHandlerFixture
	}{})
	ctrl.Setup()

	fix.ServeMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(dueoddeSearchResponse))
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

func TestOpenStreetMapQuery(t *testing.T) {
	t.Parallel()

	fix, ctrl := fixture.Init(t, &struct {
		*OpenStreetMapFixture
		*HTTPHandlerFixture
	}{})
	ctrl.Setup()

	var req *http.Request

	fix.ServeMux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) { req = r })

	_, err := fix.Search("dueodde")
	if err != nil {
		t.Errorf("Search returned an error: %v", err)
	}

	actual := req.URL.Query()["q"]
	expected := []string{"dueodde"}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Bad query. Expected %v - got %v", expected, actual)
	}
}

func TestCleanup(t *testing.T) {
	t.Parallel()

	recorder := &CleanupRecorder{TB: t}
	fix, ctrl := fixture.Init(recorder, &struct {
		*OpenStreetMapFixture
	}{})
	ctrl.Setup()

	if !fix.OpenStreetMapFixture.ServerFixture.IsOpen() {
		t.Fatal("Server should have been open")
	}

	recorder.Replay()
	if fix.OpenStreetMapFixture.ServerFixture.IsOpen() {
		t.Fatal("Server should have been closed")
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
	if !f.IsOpen() {
		f.Server = httptest.NewServer(f.Handler.ServeMux)
	}
}

func (f *HTTPServerFixture) Cleanup() {
	f.Close()
}

func (f *HTTPServerFixture) Close() {
	if f.IsOpen() {
		f.Server.Close()
		f.Server = nil
	}
}

func (f *HTTPServerFixture) IsOpen() bool {
	return f.Server != nil
}

type OpenStreetMapFixture struct {
	fixture.Fixture
	OpenStreetMap
	ServerFixture *HTTPServerFixture
}

func (f *OpenStreetMapFixture) Setup() {
	f.OpenStreetMap.BaseURL = f.ServerFixture.URL
}

type CleanupRecorder struct {
	testing.TB
	cleanups []func()
}

func (r *CleanupRecorder) Cleanup(f func()) {
	r.cleanups = append(r.cleanups, f)
	// r.TB.Cleanup(f)
}

func (r *CleanupRecorder) Replay() {
	for _, f := range r.cleanups {
		f()
	}
}

const dueoddeSearchResponse = `[{"place_id":135627303,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":2960711590,"lat":"54.9902757","lon":"15.0759924","class":"place","type":"locality","place_rank":25,"importance":0.06672011825741969,"addresstype":"locality","name":"Dueodde","display_name":"Dueodde, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9802757","55.0002757","15.0659924","15.0859924"]},{"place_id":136558536,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":792653018,"lat":"54.9935299","lon":"15.0781951","class":"highway","type":"bus_stop","place_rank":30,"importance":5.345159075298722e-05,"addresstype":"highway","name":"Dueodde","display_name":"Dueodde, Sirenevej, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9934799","54.9935799","15.0781451","15.0782451"]},{"place_id":136968598,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"node","osm_id":10006025636,"lat":"54.9954556","lon":"15.0381217","class":"tourism","type":"information","place_rank":30,"importance":9.99999999995449e-06,"addresstype":"tourism","name":"Dueodde","display_name":"Dueodde, Udegårdsvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9954056","54.9955056","15.0380717","15.0381717"]},{"place_id":134967504,"licence":"Data © OpenStreetMap contributors, ODbL 1.0. http://osm.org/copyright","osm_type":"way","osm_id":528547308,"lat":"54.991797500000004","lon":"15.074294876767677","class":"man_made","type":"lighthouse","place_rank":30,"importance":0.3797855310791376,"addresstype":"man_made","name":"Dueodde Fyr","display_name":"Dueodde Fyr, Fyrvejen, Bornholm Regional Municipality, Capital Region of Denmark, Denmark","boundingbox":["54.9917777","54.9918173","15.0742643","15.0743255"]}]`
