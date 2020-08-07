#!/bin/bash

echo 'start cron job...';

cd ~/work/src/platform-batch-game/
go run main.go

echo 'success...';