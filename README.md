# Espeak Bridge
## Requirements

 - ```eseapk-ng``` in path
 - ```sox``` in path
 - 
## Goals

The goal of this emacspeak server is to have the most straightforward and easy
setup possible.  All it requires is two things that are cross platform,
espeak-ng and sox.

## Design

Since things can get triggered very quickly in emacspeak, the goals of this
server is to avoid latency. The way we accomplish this is by having a adaptive
worker pool of espeak processes.

The reason for this to design is espeak-ng has no hooks or feedback for when 
it finishes speaking something, so what we end up doing is pushing data to it
then closing the stdin and waiting for the process to close. 

So, these espeak processes exist in 3 di

Basic design:

A series of go routines and channels to pipe data around stage by stage along
the flow

Flow is:

1. emacspeak opens espeak-bridge
2. espeak-bridge opens 2 basic processes it will use
   1. sox play listening on stdin
   2. espeak-ng with default settings going to process the next batch of text
3. espeak-bridge gets TTS commands from emacspeak
4. espeak-bridge routes them to one of three paths:
   1. tone command
      1. check if tone in internal cache, if so jump to e
      2. tone generated using sox command line in ogg format
      3. raw bytes captured by emacs-bridge
      4. raw bytes stored in internal cache with tone features
      5. raw bytes piped to open play process reading from stdin
   2. sound command
      1. check if sound is internal cache, if so jump to c.
      2. raw bytes of sound captured by emacs-bridge
      3. raw bytes piped to open play process reading from stdin
   3. tts command
      1. chunk into speed, voice, pitch, amplitude, capsmode groups
		  
		 
		 
So, that is a lot, but basically you end up with a process tree like

- ```emacs``` (running emacspeak)
   - ```espeak-bridge```
      - ```play - ``` (waiting to play ogg data)
      - ```sox ... ``` (spawned just to handle generating ogg of tone)
      - ```espeak -p10 -k2 -a 100 -v en -rate 200``` (presentingly processing
        data)
     - ```espeak -p10 -k2 -a 100 -v en -rate 220``` (waiting to process next
       chunk to be sent)
