#!/bin/bash

ROOT_DIR="$1"
RELEASE_DIR="$2"
INSTALLERS_DIR="$ROOT_DIR/scripts/installers"
IMAGES_DIR="$ROOT_DIR/docs/images"

echo ">>> Copy images..."
cp "$IMAGES_DIR/logo/win.ico" "$DEPLOY_DIR/win.ico"

echo ">>> Copy installer and uninstaller..."
SCOOP_INSTALLER="$RELEASE_DIR/$APP_NAME.json"
cp "$INSTALLERS_DIR/scoop.json" "$SCOOP_INSTALLER"
declare -A REPLACER=(
    [{APP_VERSION}]="${APP_VERSION}"
    [{APP_NAME}]="${APP_NAME}"
    [{APP_DISPLAY_NAME}]="${APP_DISPLAY_NAME}"
)
for key in "${!REPLACER[@]}"; do
    sed -i "s#$key#${REPLACER[$key]}#g" "$SCOOP_INSTALLER"
done
