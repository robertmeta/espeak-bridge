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
   1. sox play listening on stdin
   2. espeak at users base pitch/rate/voice
   3. espeak at uses multipled for letters pitch/rate/voice
3. espeak-bridge gets TTS commands from emacspeak
4. espeak-bridge routes them to one of three paths:
   1. tone command
      1. check if tone in internal cache, if so jump to 5
      2. tone generated using sox command line in ogg format
      3. raw bytes captured by emacs-bridge
      4. raw bytes stored in internal cache with tone features
      5. raw bytes piped to open play process reading from stdin
   2. sound command
      1. check if sound is internal cache, if so jump to 4.2.3
      2. raw bytes of sound captured by emacs-bridge
      3. raw bytes piped to open play process reading from stdin
   3. tts command
      1. chunk into speed, voice, pitch, amplitude, capsmode groups
      2. check interal process list to see if we have a proper peak process up
         matching that, if so jump to 4.3.6.
      3. start espeak process with correct settings
      4. add process to list of handlers
      5. if exceeding max handlers, close LRU
      6. direct tts content to correct server, which will play it with right
         settings
		 
		 
So, that is a lot, but basically you end up with a process tree like

- ```emacs``` (running emacspeak)
 - ```espeak-bridge```
  - ```play - ``` (waiting to play ogg data)
  - ```sox ... ``` (spawned just to handle generating ogg of tone)
  - ```espeak -p10 -k2 -a 100 -v en -rate 200``` (waiting to handle voice output
    in this configuration)
  - ```espeak -p10 -k2 -a 100 -v en -rate 220``` (109% faster so ready to handle
    letters for example at a 1.1 multiple)
  - ```espeak -p50 -k20 -a 100 -v en -rate 200```
  - ```...```
