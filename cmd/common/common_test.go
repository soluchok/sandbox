/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"os"
	"strconv"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"github.com/trustbloc/edge-core/pkg/log"
)

const testLogModuleName = "test"

var logger = log.New(testLogModuleName)

func TestSetLogLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		resetLoggingLevels()

		SetDefaultLogLevel(logger, "debug")

		require.Equal(t, log.DEBUG, log.GetLevel(""))
	})
	t.Run("Invalid log level", func(t *testing.T) {
		resetLoggingLevels()

		SetDefaultLogLevel(logger, "mango")

		// Should remain unchanged
		require.Equal(t, log.INFO, log.GetLevel(""))
	})
}

func TestDBParams(t *testing.T) {
	t.Run("valid params", func(t *testing.T) {
		expected := &DBParameters{
			URL:     "mem://test",
			Prefix:  "prefix",
			Timeout: 30,
		}
		setEnv(t, expected)
		defer unsetEnv(t)
		cmd := &cobra.Command{}
		Flags(cmd)
		result, err := DBParams(cmd)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("use default timeout", func(t *testing.T) {
		expected := &DBParameters{
			URL:     "mem://test",
			Prefix:  "prefix",
			Timeout: DatabaseTimeoutDefault,
		}
		setEnv(t, expected)
		defer unsetEnv(t)
		err := os.Setenv(DatabaseTimeoutEnvKey, "")
		require.NoError(t, err)
		cmd := &cobra.Command{}
		Flags(cmd)
		result, err := DBParams(cmd)
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})

	t.Run("error if url is missing", func(t *testing.T) {
		expected := &DBParameters{
			Prefix:  "prefix",
			Timeout: 30,
		}
		setEnv(t, expected)
		defer unsetEnv(t)
		cmd := &cobra.Command{}
		Flags(cmd)
		_, err := DBParams(cmd)
		require.Error(t, err)
	})

	t.Run("error if prefix is missing", func(t *testing.T) {
		expected := &DBParameters{
			URL:     "mem://test",
			Timeout: 30,
		}
		setEnv(t, expected)
		defer unsetEnv(t)
		cmd := &cobra.Command{}
		Flags(cmd)
		_, err := DBParams(cmd)
		require.Error(t, err)
	})

	t.Run("error if timeout has an invalid value", func(t *testing.T) {
		expected := &DBParameters{
			URL:    "mem://test",
			Prefix: "prefix",
		}
		setEnv(t, expected)
		defer unsetEnv(t)
		err := os.Setenv(DatabaseTimeoutEnvKey, "invalid")
		require.NoError(t, err)
		cmd := &cobra.Command{}
		Flags(cmd)
		_, err = DBParams(cmd)
		require.Error(t, err)
	})
}

func TestInitEdgeStore(t *testing.T) {
	t.Run("inits ok", func(t *testing.T) {
		s, err := InitEdgeStore(&DBParameters{
			URL:     "mem://test",
			Prefix:  "test",
			Timeout: 30,
		}, log.New("test"))
		require.NoError(t, err)
		require.NotNil(t, s)
	})

	t.Run("error if url format is invalid", func(t *testing.T) {
		_, err := InitEdgeStore(&DBParameters{
			URL:     "invalid",
			Prefix:  "test",
			Timeout: 30,
		}, log.New("test"))
		require.Error(t, err)
	})

	t.Run("error if driver is not supported", func(t *testing.T) {
		_, err := InitEdgeStore(&DBParameters{
			URL:     "unsupported://test",
			Prefix:  "test",
			Timeout: 30,
		}, log.New("test"))
		require.Error(t, err)
	})

	t.Run("error if cannot connect to store", func(t *testing.T) {
		_, err := InitEdgeStore(&DBParameters{
			URL:     "mysql://test:secret@tcp(localhost:5984)",
			Prefix:  "test",
			Timeout: 1,
		}, log.New("test"))
		require.Error(t, err)
	})
}

func resetLoggingLevels() {
	log.SetLevel("", log.INFO)
}

func setEnv(t *testing.T, values *DBParameters) {
	err := os.Setenv(DatabaseURLEnvKey, values.URL)
	require.NoError(t, err)

	err = os.Setenv(DatabasePrefixEnvKey, values.Prefix)
	require.NoError(t, err)

	err = os.Setenv(DatabaseTimeoutEnvKey, strconv.FormatUint(values.Timeout, 10))
	require.NoError(t, err)
}

func unsetEnv(t *testing.T) {
	err := os.Unsetenv(DatabaseURLEnvKey)
	require.NoError(t, err)

	err = os.Unsetenv(DatabasePrefixEnvKey)
	require.NoError(t, err)

	err = os.Unsetenv(DatabaseTimeoutEnvKey)
	require.NoError(t, err)
}
