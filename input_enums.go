package canvas3d

import "github.com/Zyko0/go-sdl3/sdl"

///	Copied and adapted/modified from github.com/Zyko0/go-sdl3.
///
///	MIT License
///
///	Copyright (c) 2025 Zyko0
///
///	Permission is hereby granted, free of charge, to any person obtaining a copy
///	of this software and associated documentation files (the "Software"), to deal
///	in the Software without restriction, including without limitation the rights
///	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
///	copies of the Software, and to permit persons to whom the Software is
///	furnished to do so, subject to the following conditions:
///
///	The above copyright notice and this permission notice shall be included in all
///	copies or substantial portions of the Software.
///
///	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
///	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
///	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
///	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
///	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
///	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
///	SOFTWARE.
///
///	Source: https://github.com/Zyko0/go-sdl3

// TODO: union type
// https://github.com/libsdl-org/SDL/blob/release-3.2.2/include/SDL3/SDL_events.h#L986
type Event = sdl.Event

// SDL_EventType - The types of events that can be delivered.
// (https://wiki.libsdl.org/SDL3/SDL_EventType)
type EventType = sdl.EventType

const (
	EVENT_FIRST                         EventType = 0   // Unused (do not remove)
	EVENT_QUIT                          EventType = 256 // User-requested quit
	EVENT_TERMINATING                   EventType = 257 // The application is being terminated by the OS. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationWillTerminate() Called on Android in onDestroy()
	EVENT_LOW_MEMORY                    EventType = 258 // The application is low on memory, free memory if possible. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationDidReceiveMemoryWarning() Called on Android in onTrimMemory()
	EVENT_WILL_ENTER_BACKGROUND         EventType = 259 // The application is about to enter the background. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationWillResignActive() Called on Android in onPause()
	EVENT_DID_ENTER_BACKGROUND          EventType = 260 // The application did enter the background and may not get CPU for some time. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationDidEnterBackground() Called on Android in onPause()
	EVENT_WILL_ENTER_FOREGROUND         EventType = 261 // The application is about to enter the foreground. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationWillEnterForeground() Called on Android in onResume()
	EVENT_DID_ENTER_FOREGROUND          EventType = 262 // The application is now interactive. This event must be handled in a callback set with SDL_AddEventWatch(). Called on iOS in applicationDidBecomeActive() Called on Android in onResume()
	EVENT_LOCALE_CHANGED                EventType = 263 // The user's locale preferences have changed.
	EVENT_SYSTEM_THEME_CHANGED          EventType = 264 // The system theme changed
	EVENT_DISPLAY_ORIENTATION           EventType = 337 // Display orientation has changed to data1
	EVENT_DISPLAY_ADDED                 EventType = 338 // Display has been added to the system
	EVENT_DISPLAY_REMOVED               EventType = 339 // Display has been removed from the system
	EVENT_DISPLAY_MOVED                 EventType = 340 // Display has changed position
	EVENT_DISPLAY_DESKTOP_MODE_CHANGED  EventType = 341 // Display has changed desktop mode
	EVENT_DISPLAY_CURRENT_MODE_CHANGED  EventType = 342 // Display has changed current mode
	EVENT_DISPLAY_CONTENT_SCALE_CHANGED EventType = 343 // Display has changed content scale
	EVENT_DISPLAY_USABLE_BOUNDS_CHANGED EventType = 344 // Display has changed usable bounds
	EVENT_DISPLAY_FIRST                 EventType = 337
	EVENT_DISPLAY_LAST                  EventType = 344
	EVENT_WINDOW_SHOWN                  EventType = 514 // Window has been shown
	EVENT_WINDOW_HIDDEN                 EventType = 515 // Window has been hidden
	EVENT_WINDOW_EXPOSED                EventType = 516 // Window has been exposed and should be redrawn, and can be redrawn directly from event watchers for this event. data1 is 1 for live-resize expose events, 0 otherwise.
	EVENT_WINDOW_MOVED                  EventType = 517 // Window has been moved to data1, data2
	EVENT_WINDOW_RESIZED                EventType = 518 // Window has been resized to data1xdata2
	EVENT_WINDOW_PIXEL_SIZE_CHANGED     EventType = 519 // The pixel size of the window has changed to data1xdata2
	EVENT_WINDOW_METAL_VIEW_RESIZED     EventType = 520 // The pixel size of a Metal view associated with the window has changed
	EVENT_WINDOW_MINIMIZED              EventType = 521 // Window has been minimized
	EVENT_WINDOW_MAXIMIZED              EventType = 522 // Window has been maximized
	EVENT_WINDOW_RESTORED               EventType = 523 // Window has been restored to normal size and position
	EVENT_WINDOW_MOUSE_ENTER            EventType = 524 // Window has gained mouse focus
	EVENT_WINDOW_MOUSE_LEAVE            EventType = 525 // Window has lost mouse focus
	EVENT_WINDOW_FOCUS_GAINED           EventType = 526 // Window has gained keyboard focus
	EVENT_WINDOW_FOCUS_LOST             EventType = 527 // Window has lost keyboard focus
	EVENT_WINDOW_CLOSE_REQUESTED        EventType = 528 // The window manager requests that the window be closed
	EVENT_WINDOW_HIT_TEST               EventType = 529 // Window had a hit test that wasn't SDL_HITTEST_NORMAL
	EVENT_WINDOW_ICCPROF_CHANGED        EventType = 530 // The window's ICC profile has changed
	EVENT_WINDOW_DISPLAY_CHANGED        EventType = 531 // Window has been moved to display data1
	EVENT_WINDOW_DISPLAY_SCALE_CHANGED  EventType = 532 // Window display scale has been changed
	EVENT_WINDOW_SAFE_AREA_CHANGED      EventType = 533 // The window safe area has been changed
	EVENT_WINDOW_OCCLUDED               EventType = 534 // The window has been occluded
	EVENT_WINDOW_ENTER_FULLSCREEN       EventType = 535 // The window has entered fullscreen mode
	EVENT_WINDOW_LEAVE_FULLSCREEN       EventType = 536 // The window has left fullscreen mode
	EVENT_WINDOW_DESTROYED              EventType = 537 // The window with the associated ID is being or has been destroyed. If this message is being handled in an event watcher, the window handle is still valid and can still be used to retrieve any properties associated with the window. Otherwise, the handle has already been destroyed and all resources associated with it are invalid
	EVENT_WINDOW_HDR_STATE_CHANGED      EventType = 538 // Window HDR properties have changed
	EVENT_WINDOW_FIRST                  EventType = 514
	EVENT_WINDOW_LAST                   EventType = 538
	EVENT_KEY_DOWN                      EventType = 768  // Key pressed
	EVENT_KEY_UP                        EventType = 769  // Key released
	EVENT_TEXT_EDITING                  EventType = 770  // Keyboard text editing (composition)
	EVENT_TEXT_INPUT                    EventType = 771  // Keyboard text input
	EVENT_KEYMAP_CHANGED                EventType = 772  // Keymap changed due to a system event such as an input language or keyboard layout change.
	EVENT_KEYBOARD_ADDED                EventType = 773  // A new keyboard has been inserted into the system
	EVENT_KEYBOARD_REMOVED              EventType = 774  // A keyboard has been removed
	EVENT_TEXT_EDITING_CANDIDATES       EventType = 775  // Keyboard text editing candidates
	EVENT_SCREEN_KEYBOARD_SHOWN         EventType = 776  // The on-screen keyboard has been shown
	EVENT_SCREEN_KEYBOARD_HIDDEN        EventType = 777  // The on-screen keyboard has been hidden
	EVENT_MOUSE_MOTION                  EventType = 1024 // Mouse moved
	EVENT_MOUSE_BUTTON_DOWN             EventType = 1025 // Mouse button pressed
	EVENT_MOUSE_BUTTON_UP               EventType = 1026 // Mouse button released
	EVENT_MOUSE_WHEEL                   EventType = 1027 // Mouse wheel motion
	EVENT_MOUSE_ADDED                   EventType = 1028 // A new mouse has been inserted into the system
	EVENT_MOUSE_REMOVED                 EventType = 1029 // A mouse has been removed
	EVENT_JOYSTICK_AXIS_MOTION          EventType = 1536 // Joystick axis motion
	EVENT_JOYSTICK_BALL_MOTION          EventType = 1537 // Joystick trackball motion
	EVENT_JOYSTICK_HAT_MOTION           EventType = 1538 // Joystick hat position change
	EVENT_JOYSTICK_BUTTON_DOWN          EventType = 1539 // Joystick button pressed
	EVENT_JOYSTICK_BUTTON_UP            EventType = 1540 // Joystick button released
	EVENT_JOYSTICK_ADDED                EventType = 1541 // A new joystick has been inserted into the system
	EVENT_JOYSTICK_REMOVED              EventType = 1542 // An opened joystick has been removed
	EVENT_JOYSTICK_BATTERY_UPDATED      EventType = 1543 // Joystick battery level change
	EVENT_JOYSTICK_UPDATE_COMPLETE      EventType = 1544 // Joystick update is complete
	EVENT_GAMEPAD_AXIS_MOTION           EventType = 1616 // Gamepad axis motion
	EVENT_GAMEPAD_BUTTON_DOWN           EventType = 1617 // Gamepad button pressed
	EVENT_GAMEPAD_BUTTON_UP             EventType = 1618 // Gamepad button released
	EVENT_GAMEPAD_ADDED                 EventType = 1619 // A new gamepad has been inserted into the system
	EVENT_GAMEPAD_REMOVED               EventType = 1620 // A gamepad has been removed
	EVENT_GAMEPAD_REMAPPED              EventType = 1621 // The gamepad mapping was updated
	EVENT_GAMEPAD_TOUCHPAD_DOWN         EventType = 1622 // Gamepad touchpad was touched
	EVENT_GAMEPAD_TOUCHPAD_MOTION       EventType = 1623 // Gamepad touchpad finger was moved
	EVENT_GAMEPAD_TOUCHPAD_UP           EventType = 1624 // Gamepad touchpad finger was lifted
	EVENT_GAMEPAD_SENSOR_UPDATE         EventType = 1625 // Gamepad sensor was updated
	EVENT_GAMEPAD_UPDATE_COMPLETE       EventType = 1626 // Gamepad update is complete
	EVENT_GAMEPAD_STEAM_HANDLE_UPDATED  EventType = 1627 // Gamepad Steam handle has changed
	EVENT_FINGER_DOWN                   EventType = 1792
	EVENT_FINGER_UP                     EventType = 1793
	EVENT_FINGER_MOTION                 EventType = 1794
	EVENT_FINGER_CANCELED               EventType = 1795
	EVENT_PINCH_BEGIN                   EventType = 1808 // Pinch gesture started
	EVENT_PINCH_UPDATE                  EventType = 1809 // Pinch gesture updated
	EVENT_PINCH_END                     EventType = 1810 // Pinch gesture ended
	EVENT_CLIPBOARD_UPDATE              EventType = 2304 // The clipboard changed
	EVENT_DROP_FILE                     EventType = 4096 // The system requests a file open
	EVENT_DROP_TEXT                     EventType = 4097 // text/plain drag-and-drop event
	EVENT_DROP_BEGIN                    EventType = 4098 // A new set of drops is beginning (NULL filename)
	EVENT_DROP_COMPLETE                 EventType = 4099 // Current set of drops is now complete (NULL filename)
	EVENT_DROP_POSITION                 EventType = 4100 // Position while moving over the window
	EVENT_AUDIO_DEVICE_ADDED            EventType = 4352 // A new audio device is available
	EVENT_AUDIO_DEVICE_REMOVED          EventType = 4353 // An audio device has been removed.
	EVENT_AUDIO_DEVICE_FORMAT_CHANGED   EventType = 4354 // An audio device's format has been changed by the system.
	EVENT_SENSOR_UPDATE                 EventType = 4608 // A sensor was updated
	EVENT_PEN_PROXIMITY_IN              EventType = 4864 // Pressure-sensitive pen has become available
	EVENT_PEN_PROXIMITY_OUT             EventType = 4865 // Pressure-sensitive pen has become unavailable
	EVENT_PEN_DOWN                      EventType = 4866 // Pressure-sensitive pen touched drawing surface
	EVENT_PEN_UP                        EventType = 4867 // Pressure-sensitive pen stopped touching drawing surface
	EVENT_PEN_BUTTON_DOWN               EventType = 4868 // Pressure-sensitive pen button pressed
	EVENT_PEN_BUTTON_UP                 EventType = 4869 // Pressure-sensitive pen button released
	EVENT_PEN_MOTION                    EventType = 4870 // Pressure-sensitive pen is moving on the tablet
	EVENT_PEN_AXIS                      EventType = 4871 // Pressure-sensitive pen angle/pressure/etc changed
	EVENT_CAMERA_DEVICE_ADDED           EventType = 5120 // A new camera device is available
	EVENT_CAMERA_DEVICE_REMOVED         EventType = 5121 // A camera device has been removed.
	EVENT_CAMERA_DEVICE_APPROVED        EventType = 5122 // A camera device has been approved for use by the user.
	EVENT_CAMERA_DEVICE_DENIED          EventType = 5123 // A camera device has been denied for use by the user.
	EVENT_RENDER_TARGETS_RESET          EventType = 8192 // The render targets have been reset and their contents need to be updated
	EVENT_RENDER_DEVICE_RESET           EventType = 8193 // The device has been reset and all textures need to be recreated
	EVENT_RENDER_DEVICE_LOST            EventType = 8194 // The device has been lost and can't be recovered.
	EVENT_PRIVATE0                      EventType = 16384
	EVENT_PRIVATE1                      EventType = 16385
	EVENT_PRIVATE2                      EventType = 16386
	EVENT_PRIVATE3                      EventType = 16387
	EVENT_POLL_SENTINEL                 EventType = 32512 // Signals the end of an event poll cycle
	EVENT_USER                          EventType = 32768
	EVENT_LAST                          EventType = 65535
	EVENT_ENUM_PADDING                  EventType = 2147483647
)

// SDL_Keycode - The SDL virtual key representation.
// (https://wiki.libsdl.org/SDL3/SDL_Keycode)
type Keycode = sdl.Keycode

const (
	K_EXTENDED_MASK = (1 << 29)
	K_SCANCODE_MASK = (1 << 30)

	K_UNKNOWN              Keycode = 0x00000000 /**< 0 */
	K_RETURN               Keycode = 0x0000000d /**< '\r' */
	K_ESCAPE               Keycode = 0x0000001b /**< '\x1B' */
	K_BACKSPACE            Keycode = 0x00000008 /**< '\b' */
	K_TAB                  Keycode = 0x00000009 /**< '\t' */
	K_SPACE                Keycode = 0x00000020 /**< ' ' */
	K_EXCLAIM              Keycode = 0x00000021 /**< '!' */
	K_DBLAPOSTROPHE        Keycode = 0x00000022 /**< '"' */
	K_HASH                 Keycode = 0x00000023 /**< '#' */
	K_DOLLAR               Keycode = 0x00000024 /**< '$' */
	K_PERCENT              Keycode = 0x00000025 /**< '%' */
	K_AMPERSAND            Keycode = 0x00000026 /**< '&' */
	K_APOSTROPHE           Keycode = 0x00000027 /**< '\'' */
	K_LEFTPAREN            Keycode = 0x00000028 /**< '(' */
	K_RIGHTPAREN           Keycode = 0x00000029 /**< ')' */
	K_ASTERISK             Keycode = 0x0000002a /**< '*' */
	K_PLUS                 Keycode = 0x0000002b /**< '+' */
	K_COMMA                Keycode = 0x0000002c /**< ',' */
	K_MINUS                Keycode = 0x0000002d /**< '-' */
	K_PERIOD               Keycode = 0x0000002e /**< '.' */
	K_SLASH                Keycode = 0x0000002f /**< '/' */
	K_0                    Keycode = 0x00000030 /**< '0' */
	K_1                    Keycode = 0x00000031 /**< '1' */
	K_2                    Keycode = 0x00000032 /**< '2' */
	K_3                    Keycode = 0x00000033 /**< '3' */
	K_4                    Keycode = 0x00000034 /**< '4' */
	K_5                    Keycode = 0x00000035 /**< '5' */
	K_6                    Keycode = 0x00000036 /**< '6' */
	K_7                    Keycode = 0x00000037 /**< '7' */
	K_8                    Keycode = 0x00000038 /**< '8' */
	K_9                    Keycode = 0x00000039 /**< '9' */
	K_COLON                Keycode = 0x0000003a /**< ':' */
	K_SEMICOLON            Keycode = 0x0000003b /**< ';' */
	K_LESS                 Keycode = 0x0000003c /**< '<' */
	K_EQUALS               Keycode = 0x0000003d /**< '=' */
	K_GREATER              Keycode = 0x0000003e /**< '>' */
	K_QUESTION             Keycode = 0x0000003f /**< '?' */
	K_AT                   Keycode = 0x00000040 /**< '@' */
	K_LEFTBRACKET          Keycode = 0x0000005b /**< '[' */
	K_BACKSLASH            Keycode = 0x0000005c /**< '\\' */
	K_RIGHTBRACKET         Keycode = 0x0000005d /**< ']' */
	K_CARET                Keycode = 0x0000005e /**< '^' */
	K_UNDERSCORE           Keycode = 0x0000005f /**< '_' */
	K_GRAVE                Keycode = 0x00000060 /**< '`' */
	K_A                    Keycode = 0x00000061 /**< 'a' */
	K_B                    Keycode = 0x00000062 /**< 'b' */
	K_C                    Keycode = 0x00000063 /**< 'c' */
	K_D                    Keycode = 0x00000064 /**< 'd' */
	K_E                    Keycode = 0x00000065 /**< 'e' */
	K_F                    Keycode = 0x00000066 /**< 'f' */
	K_G                    Keycode = 0x00000067 /**< 'g' */
	K_H                    Keycode = 0x00000068 /**< 'h' */
	K_I                    Keycode = 0x00000069 /**< 'i' */
	K_J                    Keycode = 0x0000006a /**< 'j' */
	K_K                    Keycode = 0x0000006b /**< 'k' */
	K_L                    Keycode = 0x0000006c /**< 'l' */
	K_M                    Keycode = 0x0000006d /**< 'm' */
	K_N                    Keycode = 0x0000006e /**< 'n' */
	K_O                    Keycode = 0x0000006f /**< 'o' */
	K_P                    Keycode = 0x00000070 /**< 'p' */
	K_Q                    Keycode = 0x00000071 /**< 'q' */
	K_R                    Keycode = 0x00000072 /**< 'r' */
	K_S                    Keycode = 0x00000073 /**< 's' */
	K_T                    Keycode = 0x00000074 /**< 't' */
	K_U                    Keycode = 0x00000075 /**< 'u' */
	K_V                    Keycode = 0x00000076 /**< 'v' */
	K_W                    Keycode = 0x00000077 /**< 'w' */
	K_X                    Keycode = 0x00000078 /**< 'x' */
	K_Y                    Keycode = 0x00000079 /**< 'y' */
	K_Z                    Keycode = 0x0000007a /**< 'z' */
	K_LEFTBRACE            Keycode = 0x0000007b /**< '{' */
	K_PIPE                 Keycode = 0x0000007c /**< '|' */
	K_RIGHTBRACE           Keycode = 0x0000007d /**< '}' */
	K_TILDE                Keycode = 0x0000007e /**< '~' */
	K_DELETE               Keycode = 0x0000007f /**< '\x7F' */
	K_PLUSMINUS            Keycode = 0x000000b1 /**< '\xB1' */
	K_CAPSLOCK             Keycode = 0x40000039 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CAPSLOCK) */
	K_F1                   Keycode = 0x4000003a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F1) */
	K_F2                   Keycode = 0x4000003b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F2) */
	K_F3                   Keycode = 0x4000003c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F3) */
	K_F4                   Keycode = 0x4000003d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F4) */
	K_F5                   Keycode = 0x4000003e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F5) */
	K_F6                   Keycode = 0x4000003f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F6) */
	K_F7                   Keycode = 0x40000040 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F7) */
	K_F8                   Keycode = 0x40000041 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F8) */
	K_F9                   Keycode = 0x40000042 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F9) */
	K_F10                  Keycode = 0x40000043 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F10) */
	K_F11                  Keycode = 0x40000044 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F11) */
	K_F12                  Keycode = 0x40000045 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F12) */
	K_PRINTSCREEN          Keycode = 0x40000046 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PRINTSCREEN) */
	K_SCROLLLOCK           Keycode = 0x40000047 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SCROLLLOCK) */
	K_PAUSE                Keycode = 0x40000048 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PAUSE) */
	K_INSERT               Keycode = 0x40000049 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_INSERT) */
	K_HOME                 Keycode = 0x4000004a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_HOME) */
	K_PAGEUP               Keycode = 0x4000004b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PAGEUP) */
	K_END                  Keycode = 0x4000004d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_END) */
	K_PAGEDOWN             Keycode = 0x4000004e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PAGEDOWN) */
	K_RIGHT                Keycode = 0x4000004f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RIGHT) */
	K_LEFT                 Keycode = 0x40000050 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_LEFT) */
	K_DOWN                 Keycode = 0x40000051 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_DOWN) */
	K_UP                   Keycode = 0x40000052 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_UP) */
	K_NUMLOCKCLEAR         Keycode = 0x40000053 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_NUMLOCKCLEAR) */
	K_KP_DIVIDE            Keycode = 0x40000054 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_DIVIDE) */
	K_KP_MULTIPLY          Keycode = 0x40000055 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MULTIPLY) */
	K_KP_MINUS             Keycode = 0x40000056 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MINUS) */
	K_KP_PLUS              Keycode = 0x40000057 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_PLUS) */
	K_KP_ENTER             Keycode = 0x40000058 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_ENTER) */
	K_KP_1                 Keycode = 0x40000059 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_1) */
	K_KP_2                 Keycode = 0x4000005a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_2) */
	K_KP_3                 Keycode = 0x4000005b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_3) */
	K_KP_4                 Keycode = 0x4000005c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_4) */
	K_KP_5                 Keycode = 0x4000005d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_5) */
	K_KP_6                 Keycode = 0x4000005e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_6) */
	K_KP_7                 Keycode = 0x4000005f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_7) */
	K_KP_8                 Keycode = 0x40000060 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_8) */
	K_KP_9                 Keycode = 0x40000061 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_9) */
	K_KP_0                 Keycode = 0x40000062 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_0) */
	K_KP_PERIOD            Keycode = 0x40000063 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_PERIOD) */
	K_APPLICATION          Keycode = 0x40000065 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_APPLICATION) */
	K_POWER                Keycode = 0x40000066 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_POWER) */
	K_KP_EQUALS            Keycode = 0x40000067 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_EQUALS) */
	K_F13                  Keycode = 0x40000068 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F13) */
	K_F14                  Keycode = 0x40000069 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F14) */
	K_F15                  Keycode = 0x4000006a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F15) */
	K_F16                  Keycode = 0x4000006b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F16) */
	K_F17                  Keycode = 0x4000006c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F17) */
	K_F18                  Keycode = 0x4000006d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F18) */
	K_F19                  Keycode = 0x4000006e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F19) */
	K_F20                  Keycode = 0x4000006f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F20) */
	K_F21                  Keycode = 0x40000070 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F21) */
	K_F22                  Keycode = 0x40000071 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F22) */
	K_F23                  Keycode = 0x40000072 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F23) */
	K_F24                  Keycode = 0x40000073 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_F24) */
	K_EXECUTE              Keycode = 0x40000074 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_EXECUTE) */
	K_HELP                 Keycode = 0x40000075 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_HELP) */
	K_MENU                 Keycode = 0x40000076 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MENU) */
	K_SELECT               Keycode = 0x40000077 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SELECT) */
	K_STOP                 Keycode = 0x40000078 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_STOP) */
	K_AGAIN                Keycode = 0x40000079 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AGAIN) */
	K_UNDO                 Keycode = 0x4000007a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_UNDO) */
	K_CUT                  Keycode = 0x4000007b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CUT) */
	K_COPY                 Keycode = 0x4000007c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_COPY) */
	K_PASTE                Keycode = 0x4000007d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PASTE) */
	K_FIND                 Keycode = 0x4000007e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_FIND) */
	K_MUTE                 Keycode = 0x4000007f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MUTE) */
	K_VOLUMEUP             Keycode = 0x40000080 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_VOLUMEUP) */
	K_VOLUMEDOWN           Keycode = 0x40000081 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_VOLUMEDOWN) */
	K_KP_COMMA             Keycode = 0x40000085 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_COMMA) */
	K_KP_EQUALSAS400       Keycode = 0x40000086 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_EQUALSAS400) */
	K_ALTERASE             Keycode = 0x40000099 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_ALTERASE) */
	K_SYSREQ               Keycode = 0x4000009a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SYSREQ) */
	K_CANCEL               Keycode = 0x4000009b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CANCEL) */
	K_CLEAR                Keycode = 0x4000009c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CLEAR) */
	K_PRIOR                Keycode = 0x4000009d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_PRIOR) */
	K_RETURN2              Keycode = 0x4000009e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RETURN2) */
	K_SEPARATOR            Keycode = 0x4000009f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SEPARATOR) */
	K_OUT                  Keycode = 0x400000a0 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_OUT) */
	K_OPER                 Keycode = 0x400000a1 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_OPER) */
	K_CLEARAGAIN           Keycode = 0x400000a2 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CLEARAGAIN) */
	K_CRSEL                Keycode = 0x400000a3 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CRSEL) */
	K_EXSEL                Keycode = 0x400000a4 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_EXSEL) */
	K_KP_00                Keycode = 0x400000b0 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_00) */
	K_KP_000               Keycode = 0x400000b1 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_000) */
	K_THOUSANDSSEPARATOR   Keycode = 0x400000b2 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_THOUSANDSSEPARATOR) */
	K_DECIMALSEPARATOR     Keycode = 0x400000b3 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_DECIMALSEPARATOR) */
	K_CURRENCYUNIT         Keycode = 0x400000b4 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CURRENCYUNIT) */
	K_CURRENCYSUBUNIT      Keycode = 0x400000b5 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CURRENCYSUBUNIT) */
	K_KP_LEFTPAREN         Keycode = 0x400000b6 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_LEFTPAREN) */
	K_KP_RIGHTPAREN        Keycode = 0x400000b7 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_RIGHTPAREN) */
	K_KP_LEFTBRACE         Keycode = 0x400000b8 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_LEFTBRACE) */
	K_KP_RIGHTBRACE        Keycode = 0x400000b9 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_RIGHTBRACE) */
	K_KP_TAB               Keycode = 0x400000ba /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_TAB) */
	K_KP_BACKSPACE         Keycode = 0x400000bb /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_BACKSPACE) */
	K_KP_A                 Keycode = 0x400000bc /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_A) */
	K_KP_B                 Keycode = 0x400000bd /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_B) */
	K_KP_C                 Keycode = 0x400000be /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_C) */
	K_KP_D                 Keycode = 0x400000bf /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_D) */
	K_KP_E                 Keycode = 0x400000c0 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_E) */
	K_KP_F                 Keycode = 0x400000c1 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_F) */
	K_KP_XOR               Keycode = 0x400000c2 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_XOR) */
	K_KP_POWER             Keycode = 0x400000c3 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_POWER) */
	K_KP_PERCENT           Keycode = 0x400000c4 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_PERCENT) */
	K_KP_LESS              Keycode = 0x400000c5 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_LESS) */
	K_KP_GREATER           Keycode = 0x400000c6 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_GREATER) */
	K_KP_AMPERSAND         Keycode = 0x400000c7 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_AMPERSAND) */
	K_KP_DBLAMPERSAND      Keycode = 0x400000c8 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_DBLAMPERSAND) */
	K_KP_VERTICALBAR       Keycode = 0x400000c9 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_VERTICALBAR) */
	K_KP_DBLVERTICALBAR    Keycode = 0x400000ca /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_DBLVERTICALBAR) */
	K_KP_COLON             Keycode = 0x400000cb /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_COLON) */
	K_KP_HASH              Keycode = 0x400000cc /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_HASH) */
	K_KP_SPACE             Keycode = 0x400000cd /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_SPACE) */
	K_KP_AT                Keycode = 0x400000ce /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_AT) */
	K_KP_EXCLAM            Keycode = 0x400000cf /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_EXCLAM) */
	K_KP_MEMSTORE          Keycode = 0x400000d0 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMSTORE) */
	K_KP_MEMRECALL         Keycode = 0x400000d1 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMRECALL) */
	K_KP_MEMCLEAR          Keycode = 0x400000d2 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMCLEAR) */
	K_KP_MEMADD            Keycode = 0x400000d3 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMADD) */
	K_KP_MEMSUBTRACT       Keycode = 0x400000d4 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMSUBTRACT) */
	K_KP_MEMMULTIPLY       Keycode = 0x400000d5 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMMULTIPLY) */
	K_KP_MEMDIVIDE         Keycode = 0x400000d6 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_MEMDIVIDE) */
	K_KP_PLUSMINUS         Keycode = 0x400000d7 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_PLUSMINUS) */
	K_KP_CLEAR             Keycode = 0x400000d8 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_CLEAR) */
	K_KP_CLEARENTRY        Keycode = 0x400000d9 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_CLEARENTRY) */
	K_KP_BINARY            Keycode = 0x400000da /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_BINARY) */
	K_KP_OCTAL             Keycode = 0x400000db /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_OCTAL) */
	K_KP_DECIMAL           Keycode = 0x400000dc /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_DECIMAL) */
	K_KP_HEXADECIMAL       Keycode = 0x400000dd /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_KP_HEXADECIMAL) */
	K_LCTRL                Keycode = 0x400000e0 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_LCTRL) */
	K_LSHIFT               Keycode = 0x400000e1 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_LSHIFT) */
	K_LALT                 Keycode = 0x400000e2 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_LALT) */
	K_LGUI                 Keycode = 0x400000e3 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_LGUI) */
	K_RCTRL                Keycode = 0x400000e4 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RCTRL) */
	K_RSHIFT               Keycode = 0x400000e5 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RSHIFT) */
	K_RALT                 Keycode = 0x400000e6 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RALT) */
	K_RGUI                 Keycode = 0x400000e7 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_RGUI) */
	K_MODE                 Keycode = 0x40000101 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MODE) */
	K_SLEEP                Keycode = 0x40000102 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SLEEP) */
	K_WAKE                 Keycode = 0x40000103 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_WAKE) */
	K_CHANNEL_INCREMENT    Keycode = 0x40000104 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CHANNEL_INCREMENT) */
	K_CHANNEL_DECREMENT    Keycode = 0x40000105 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CHANNEL_DECREMENT) */
	K_MEDIA_PLAY           Keycode = 0x40000106 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_PLAY) */
	K_MEDIA_PAUSE          Keycode = 0x40000107 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_PAUSE) */
	K_MEDIA_RECORD         Keycode = 0x40000108 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_RECORD) */
	K_MEDIA_FAST_FORWARD   Keycode = 0x40000109 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_FAST_FORWARD) */
	K_MEDIA_REWIND         Keycode = 0x4000010a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_REWIND) */
	K_MEDIA_NEXT_TRACK     Keycode = 0x4000010b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_NEXT_TRACK) */
	K_MEDIA_PREVIOUS_TRACK Keycode = 0x4000010c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_PREVIOUS_TRACK) */
	K_MEDIA_STOP           Keycode = 0x4000010d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_STOP) */
	K_MEDIA_EJECT          Keycode = 0x4000010e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_EJECT) */
	K_MEDIA_PLAY_PAUSE     Keycode = 0x4000010f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_PLAY_PAUSE) */
	K_MEDIA_SELECT         Keycode = 0x40000110 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_MEDIA_SELECT) */
	K_AC_NEW               Keycode = 0x40000111 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_NEW) */
	K_AC_OPEN              Keycode = 0x40000112 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_OPEN) */
	K_AC_CLOSE             Keycode = 0x40000113 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_CLOSE) */
	K_AC_EXIT              Keycode = 0x40000114 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_EXIT) */
	K_AC_SAVE              Keycode = 0x40000115 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_SAVE) */
	K_AC_PRINT             Keycode = 0x40000116 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_PRINT) */
	K_AC_PROPERTIES        Keycode = 0x40000117 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_PROPERTIES) */
	K_AC_SEARCH            Keycode = 0x40000118 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_SEARCH) */
	K_AC_HOME              Keycode = 0x40000119 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_HOME) */
	K_AC_BACK              Keycode = 0x4000011a /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_BACK) */
	K_AC_FORWARD           Keycode = 0x4000011b /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_FORWARD) */
	K_AC_STOP              Keycode = 0x4000011c /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_STOP) */
	K_AC_REFRESH           Keycode = 0x4000011d /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_REFRESH) */
	K_AC_BOOKMARKS         Keycode = 0x4000011e /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_AC_BOOKMARKS) */
	K_SOFTLEFT             Keycode = 0x4000011f /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SOFTLEFT) */
	K_SOFTRIGHT            Keycode = 0x40000120 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_SOFTRIGHT) */
	K_CALL                 Keycode = 0x40000121 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_CALL) */
	K_ENDCALL              Keycode = 0x40000122 /**< SDL_SCANCODE_TO_KEYCODE(SDL_SCANCODE_ENDCALL) */
	K_LEFT_TAB             Keycode = 0x20000001 /**< Extended key Left Tab */
	K_LEVEL5_SHIFT         Keycode = 0x20000002 /**< Extended key Level 5 Shift */
	K_MULTI_KEY_COMPOSE    Keycode = 0x20000003 /**< Extended key Multi-key Compose */
	K_LMETA                Keycode = 0x20000004 /**< Extended key Left Meta */
	K_RMETA                Keycode = 0x20000005 /**< Extended key Right Meta */
	K_LHYPER               Keycode = 0x20000006 /**< Extended key Left Hyper */
	K_RHYPER               Keycode = 0x20000007 /**< Extended key Right Hyper */
)

// SDL_Keymod - Valid key modifiers (possibly OR'd together).
// (https://wiki.libsdl.org/SDL3/SDL_Keymod)
type Keymod = sdl.Keymod

const (
	KMOD_NONE   Keymod = 0x0000                      /**< no modifier is applicable. */
	KMOD_LSHIFT Keymod = 0x0001                      /**< the left Shift key is down. */
	KMOD_RSHIFT Keymod = 0x0002                      /**< the right Shift key is down. */
	KMOD_LEVEL5 Keymod = 0x0004                      /**< the Level 5 Shift key is down. */
	KMOD_LCTRL  Keymod = 0x0040                      /**< the left Ctrl (Control) key is down. */
	KMOD_RCTRL  Keymod = 0x0080                      /**< the right Ctrl (Control) key is down. */
	KMOD_LALT   Keymod = 0x0100                      /**< the left Alt key is down. */
	KMOD_RALT   Keymod = 0x0200                      /**< the right Alt key is down. */
	KMOD_LGUI   Keymod = 0x0400                      /**< the left GUI key (often the Windows key) is down. */
	KMOD_RGUI   Keymod = 0x0800                      /**< the right GUI key (often the Windows key) is down. */
	KMOD_NUM    Keymod = 0x1000                      /**< the Num Lock key (may be located on an extended keypad) is down. */
	KMOD_CAPS   Keymod = 0x2000                      /**< the Caps Lock key is down. */
	KMOD_MODE   Keymod = 0x4000                      /**< the !AltGr key is down. */
	KMOD_SCROLL Keymod = 0x8000                      /**< the Scroll Lock key is down. */
	KMOD_CTRL   Keymod = (KMOD_LCTRL | KMOD_RCTRL)   /**< Any Ctrl key is down. */
	KMOD_SHIFT  Keymod = (KMOD_LSHIFT | KMOD_RSHIFT) /**< Any Shift key is down. */
	KMOD_ALT    Keymod = (KMOD_LALT | KMOD_RALT)     /**< Any Alt key is down. */
	KMOD_GUI    Keymod = (KMOD_LGUI | KMOD_RGUI)     /**< Any GUI key is down. */
)
