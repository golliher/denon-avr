
# denon-avr

Denon-avr is a command line utility to control a Denon AVR over TCP via the published Denon API.     

I have the AVR-X4000 which is what I developed and tested against.   Development and testing done from a computer running Apple OS X Mavericks.

The documentation used was AVRX4000_PROTOCOL(10.3.0)_V01.pdf and I imagine you can find a copy via a google search.  I've had the file long enough I don't recal where I found it.

## the itch (i.e. my motivation)

I am aware of Brad Fitzpatrick's pagage at code.google.com/p/go-avr/avr.   I read it and it helped me.  My code helped me scratch my own itch though.

Namely these three things:

1. I was interested in what the AVR returns, especially for HD Radio metadata via the HD? commmand of the Denon protocol.
2. It was a good side project to help learn how to develop in Go
3. As a go newbie, Brad's gode is a bit over my head and I couldn't figure out how to get it to print the results the AVR returns.

## usage

Run the program and pass it a Denon command.  The AVR executes the command and then, depending on the situation, will print a result.


## examples


// This command issues the "SI?" command will returns the currently selected input.  Note commands must be upper case or the will silently fail.

$ denon-avr SI?
received:  SIGAME
received:  PSRSTR OFF
received:  SVOFF

// Sets the current input to HDRADIO.  

$ denon-avr SIHDRADIO
$

// Get information about the HDRADIO current state

$ denon-avr HD?
received:  HDSIG LEV 6
received:  HDST NAME WUBL-FM 
received:  HDMLT CURRCH 1
received:  HDMLT CAST CH 2
received:  HDPTY Country           
received:  HDTITLE I Don't Dance                           
received:  HDARTIST Lee Brice  

// un-mutes

$ denon-avr MUOFF
received:  MUOFF

// mutes

$ denon-avr MUON
received:  MUON

// issue the same command again and it will be silent, presumable the AVR only returns information if there was a state change

$ denon-avr MUON
$

// Checks to see if we're muted of not

$  MU?
received:  MUON

## author's request

If you find this program interesting or useful in any way, please send me a message and let me know.
