package env

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestDefaultString(t *testing.T) {
	v := defaultString("a", "test")
	if v != "test" {
		t.Fatal("v must be test")
	}

	if err := os.Setenv("a", "test_1"); err != nil {
		t.Fatal(err)
	}

	v = defaultString("a", "test")
	if v != "test_1" {
		t.Fatal("v must be test_1")
	}
}

func TestEnv(t *testing.T) {
	cases := []struct {
		flag string
		env  string
		def  string
		val  *string
	}{
		{"region", "REGION", _defaultRegion, &Region},
		{"zone", "ZONE", _defaultZone, &Zone},
		{"deploy.env", "DEPLOY_ENV", _defaultDeployEnv, &DeployEnv},
		{"appid", "APP_ID", "", &AppID},
		{"appname", "APP_Name", "", &AppName},
		{"http.port", "DISCOVERY_HTTP_PORT", _defaultHTTPPort, &HTTPPort},
	}
	for _, tc := range cases {
		// flag set
		t.Run(fmt.Sprintf("%s: flag set", tc.env), func(t *testing.T) {
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{fmt.Sprintf("-%s=%s", tc.flag, "test")})
			if err != nil {
				t.Fatal(err)
			}
			if *tc.val != "test" {
				t.Fatal("val must be test")
			}
		})

		// flag not set, env set
		t.Run(fmt.Sprintf("%s: flag not set, env set", tc.env), func(t *testing.T) {
			*tc.val = ""
			os.Setenv(tc.env, "test2")
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{})
			if err != nil {
				t.Fatal(err)
			}
			if *tc.val != "test2" {
				t.Fatal("val must be test2")
			}
		})

		// flag not set, env not set
		t.Run(fmt.Sprintf("%s: flag not set, env not set", tc.env), func(t *testing.T) {
			*tc.val = ""
			os.Setenv(tc.env, "")
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			addFlag(fs)
			err := fs.Parse([]string{})
			if err != nil {
				t.Fatal(err)
			}
			if *tc.val != tc.def {
				t.Fatal("val must be test")
			}
		})
	}
}
