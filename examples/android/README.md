### Android example

To compile example to shared library you will need [Android NDK](https://developer.android.com/ndk/downloads/index.html).
To build Android apk you will need [Android SDK](http://developer.android.com/sdk/index.html#Other).

Export path to Android NDK, point to location where you have unpacked archive:

    export ANDROID_NDK_HOME=/opt/android-ndk

Add toolchain bin directory to PATH:

    export PATH=${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin:${PATH}

Export sysroot path:
    
    export ANDROID_SYSROOT=${ANDROID_NDK_HOME}/platforms/android-16/arch-arm

And compile shared library:

    CC=arm-linux-androideabi-gcc \
    CGO_CFLAGS="-D__ANDROID_API__=16 -I${ANDROID_NDK_HOME}/sysroot/usr/include -I${ANDROID_NDK_HOME}/sysroot/usr/include/arm-linux-androideabi --sysroot=${ANDROID_SYSROOT}" \
    CGO_LDFLAGS="-L${ANDROID_NDK_HOME}/sysroot/usr/lib -L${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/ --sysroot=${ANDROID_SYSROOT}" \
    CGO_ENABLED=1 GOOS=android GOARCH=arm \
    go build -tags static -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" -o=android/libs/armeabi-v7a/libexample.so

To build apk export path to Android SDK, point to location where you unpacked archive:

    export ANDROID_HOME=/opt/android-sdk

And build apk:

    ./gradlew assembleDebug

If everything is successfully built apk can be found in android/build/outputs/apk/debug directory.

## Windows

To do the same on a freshly installed Windows with [Go](https://golang.org/dl), [Git Bash](https://git-scm.com), [Android SDK](https://developer.android.com/sdk/index.html) installed and [Android NDK](https://developer.android.com/ndk/downloads/index.html) extracted (to keep it simple, I extracted it to the Downloads directory), we will then set the environment variables. You should be able to get to the UI by pressing Windows key and search `environment`. An item called `Edit the system environment variables` should come up and after clicking it, click `Environment Variables...`. On that UI, set the following environment variables under the `User variables for [user]` section. Please modify `lilis` to the user used on your Windows machine!


```
ANDROID_HOME=C:\Users\lilis\AppData\Local\Android\Sdk
ANDROID_NDK_HOME=C:\Users\lilis\Downloads\android\ndk-r21-windows-x86_64\android-ndk-r21
ANDROID_SYSROOT=%ANDROID_NDK_HOME%\platforms\android-16\arch-arm
JAVA_HOME=C:\Program Files\Android\Android Studio\jre
Path=%ANDROID_NDK_HOME%\toolchains\llvm\prebuilt\windows-x86_64\bin;%ANDROID_NDK_HOME%\toolchains\arm-linux-androideabi-4.9\prebuilt\windows-x86_64\bin;%Path%;%USERPROFILE%\go\bin;%USERPROFILE%\AppData\Local\Microsoft\WindowsApps;
```

For the `Path` environment variable, only the first two entries are particularly relevant to this example. We are also going to just use the Java provided by Android Studio.

After that, launch `Git Bash` and run `go get -v github.com/veandco/go-sdl2/sdl` and `cd go/src/github.com/veandco/go-sdl2/.go-sdl2-examples/examples/android`. Then you can compile the SDL2 program into a shared library by running similar command as the one for Linux:

```
CC=armv7a-linux-androideabi16-clang \
CGO_CFLAGS="-D__ANDROID_API__=16 -I${ANDROID_NDK_HOME}/sysroot/usr/include -I${ANDROID_NDK_HOME}/sysroot/usr/include/arm-linux-androideabi --sysroot=${ANDROID_SYSROOT}" \
CGO_LDFLAGS="-L${ANDROID_NDK_HOME}/sysroot/usr/lib -L${ANDROID_NDK_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/windows-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/ --sysroot=${ANDROID_SYSROOT}" \
CGO_ENABLED=1 GOOS=android GOARCH=arm \
go build -tags static -buildmode=c-shared -ldflags="-s -w -extldflags=-Wl,-soname,libexample.so" -o=android/libs/armeabi-v7a/libexample.so
```

If it complains about SDK license haven't been accepted, launch `Android Studio` and click on `Configure > SDK Manager` on the bottom right of the dialog. Then under `System settings > Android SDK`, install the SDK and accept license for the API version that matches the one on the error message (e.g. Android 9.0 (Pie) which has API level 28 and Android 4.1 (Jelly Bean) which has API level 16).

After that, try running the command above again. If no more error messages show up then you should be able to run `./gradlew assembleDebug`! The resulting APK will be at `android/build/outputs/apk/debug` directory.
