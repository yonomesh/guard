package unicmd

import (
	"fmt"
	"runtime/debug"

	"uni"
)

type moduleInfo struct {
	guardModuleID string
	golangModule  *debug.Module
	err           error
}

func cmdStart(fl Flags) (int, error) {
	configFlag := fl.String("config")
	pidfileFlag := fl.String("pidfile")
	profileFlag := fl.String("profile")
	fmt.Println("cmdStart testing func")

	fmt.Println(configFlag)
	fmt.Println(pidfileFlag)
	fmt.Println(profileFlag)

	return uni.ExitCodeSuccess, nil
}

func cmdTest(fl Flags) (int, error) {
	helloFlag := fl.String("hello")
	fmt.Println("test cmd")
	fmt.Println(helloFlag)
	return uni.ExitCodeSuccess, nil
}

func cmdRun(fl Flags) (int, error) {
	uni.TrapSignals()

	// set up buffered logging for early startup
	// so that we can hold onto logs until after
	// the config is loaded (or fails to load)
	// so that we can write the logs to the user's
	// configured output. we must be sure to flush
	// on any error before the config is loaded.
	logger, defaultLogger, logBuffer := uni.BufferedLog()

	undoMaxProcs := setResourceLimits(logger)
	defer undoMaxProcs()
	// release the local reference to the undo function so it can be GC'd;
	// the deferred call above has already captured the actual function value.
	undoMaxProcs = nil //nolint:ineffassign,wastedassign

	profileFlag := fl.String("profile")
	configFlag := fl.String("config")
	resumeFlag := fl.Bool("resume")
	pidfileFlag := fl.String("pidfile")

	// load all additional envs as soon as possible
	err := handleEnvFileFlag(fl)
	if err != nil {
		logBuffer.FlushTo(defaultLogger)
		return uni.ExitCodeFailedStartup, err
	}

	// TODO
	// if we are supposed to print the environment, do that first
	// if printEnvFlag {
	// 	printEnvironment()
	// }

	// TODO
	// // load the config, depending on flags
	var config []byte
	// if resumeFlag {
	// 	config, err = os.ReadFile(caddy.ConfigAutosavePath)
	// 	if errors.Is(err, fs.ErrNotExist) {
	// 		// not a bad error; just can't resume if autosave file doesn't exist
	// 		logger.Info("no autosave file exists", zap.String("autosave_file", caddy.ConfigAutosavePath))
	// 		resumeFlag = false
	// 	} else if err != nil {
	// 		logBuffer.FlushTo(defaultLogger)
	// 		return caddy.ExitCodeFailedStartup, err
	// 	} else {
	// 		if configFlag == "" {
	// 			logger.Info("resuming from last configuration",
	// 				zap.String("autosave_file", caddy.ConfigAutosavePath))
	// 		} else {
	// 			// if they also specified a config file, user should be aware that we're not
	// 			// using it (doing so could lead to data/config loss by overwriting!)
	// 			logger.Warn("--config and --resume flags were used together; ignoring --config and resuming from last configuration",
	// 				zap.String("autosave_file", caddy.ConfigAutosavePath))
	// 		}
	// 	}
	// }

	// we don't use 'else' here since this value might have been changed in 'if' block; i.e. not mutually exclusive
	var configFile string
	// var adapterUsed string
	if !resumeFlag {
		config, configFile, err = LoadConfig(configFlag)
		if err != nil {
			logBuffer.FlushTo(defaultLogger)
			return uni.ExitCodeFailedStartup, err
		}
	}

	// TODO
	// create pidfile now, in case loading config takes a while (issue #5477)
	if pidfileFlag != "" {
		// err := caddy.PIDFile(pidfileFlag)
		// if err != nil {
		// 	logger.Error("unable to write PID file",
		// 		zap.String("pidfile", pidfileFlag),
		// 		zap.Error(err))
		// }
	}

	// TODO
	// If we have a source config file (we're running via 'caddy run --config ...'),
	// record it so SIGUSR1 can reload from the same file. Also provide a callback
	// that knows how to load/adapt that source when requested by the main process.
	if configFile != "" {
		// caddy.SetLastConfig(configFile, adapterUsed, func(file, adapter string) error {
		// 	cfg, _, _, err := LoadConfig(file, adapter)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	return caddy.Load(cfg, true)
		// })
	}

	// run the initial config

	_ = pidfileFlag
	_ = config
	_ = profileFlag

	select {}
}
