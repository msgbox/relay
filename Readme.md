MsgBox Relay
================

This is the main relay service for a MsgBox. It's written in Go and uses an AMQP queue to
hold incoming and outgoing messages.

It uses Protocol Buffers for storing and passing around data.

The [MsgBox Spec](https://github.com/msgbox/Spec) is still a work in progress as is
this relay system.

It's more of a way for me to scratch an itch that I have with the way email is used and
a project to use as a way to learn how Go works.

# Install Instructions

To-Do