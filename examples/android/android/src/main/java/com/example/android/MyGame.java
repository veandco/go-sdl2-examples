package com.example.android;

import org.libsdl.app.SDLActivity;

public class MyGame extends SDLActivity {

    @Override
    protected String[] getLibraries() {
        return new String[] {
            "example"
        };
    }

    @Override
    protected String getMainFunction() {
        return "SDL_main";
    }

}
