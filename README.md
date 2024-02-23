# Espeak Bridge

## Goals

The goal of this emacspeak server is to have the most straightforward and easy
setup possible.  All it requires is two things that are cross platform, espeak
and sox.

## Design

Since things can get triggered very quickly in emacspeak, the goals of this
server is to avoid overwhelming the system with process spawning.  Processes on
Windows are very expensive, and on Mac have extra overhead due to the way the
security system functions.

Basic design:

A series of go routines and channels to pipe data around stage by stage along
the flow

Flow is: 

1. emacspeak opens espeak-bridge 
2. espeak-bridge opens 3 basic processes it will use
2.1. sox play listening on stdin
2.2. espeak at users base pitch/rate/voice
2.3. espeak at uses multipled for letters pitch/rate/voice
3. espeak-bridge gets TTS commands from emacspeak 
4. espeak-bridge routes them to one of three paths:
4.1. tone command
4.1.1 check if tone in internal cache, if so jump to 4.1.5
4.1.2. tone generated using sox command line in ogg format
4.1.3. raw bytes captured by emacs-bridge
4.1.4. raw bytes stored in internal cache with tone features
4.1.5. raw bytes piped to open play process reading from stdin
4.2. sound command
4.2.1. check if sound is internal cache, if so jump to 4.2.3
4.2.2. raw bytes of sound captured by emacs-bridge
4.2.3. raw bytes piped to open play process reading from stdin
4.3. tts command
4.3.1. chunk into speed, voice, pitch, amplitude, capsmode groups
4.3.2. check interal process list to see if we have a proper 
       espeak process up matching that, if so jump to 4.3.6.
4.3.3. start espeak process with correct settings 
4.3.4. add process to list of handlers
4.3.5. if exceeding max handlers, close LRU
4.3.6. direct tts content to correct server, which will play it with right settings
