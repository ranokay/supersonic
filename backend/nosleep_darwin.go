//go:build darwin

package backend

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit
#include <CoreFoundation/CoreFoundation.h>
#include <IOKit/pwr_mgt/IOPMLib.h>
#include <pthread.h>

static IOPMAssertionID supersonicSleepAssertion = kIOPMNullAssertionID;
static pthread_mutex_t supersonicSleepAssertionMu = PTHREAD_MUTEX_INITIALIZER;

static void supersonic_set_system_sleep_disabled(int disable) {
	pthread_mutex_lock(&supersonicSleepAssertionMu);
	if (disable) {
		if (supersonicSleepAssertion != kIOPMNullAssertionID) {
			pthread_mutex_unlock(&supersonicSleepAssertionMu);
			return;
		}
		CFStringRef reason = CFSTR("Supersonic playback");
		IOPMAssertionCreateWithName(kIOPMAssertionTypeNoIdleSleep,
		                            kIOPMAssertionLevelOn,
		                            reason,
		                            &supersonicSleepAssertion);
		pthread_mutex_unlock(&supersonicSleepAssertionMu);
		return;
	}
	if (supersonicSleepAssertion != kIOPMNullAssertionID) {
		IOPMAssertionRelease(supersonicSleepAssertion);
		supersonicSleepAssertion = kIOPMNullAssertionID;
	}
	pthread_mutex_unlock(&supersonicSleepAssertionMu);
}
*/
import "C"

func SetSystemSleepDisabled(disable bool) {
	if disable {
		C.supersonic_set_system_sleep_disabled(1)
		return
	}
	C.supersonic_set_system_sleep_disabled(0)
}
