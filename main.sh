if [ ! -n "$SHOULD_START_WORKFLOWS" ]; then
　　echo "$SHOULD_START_WORKFLOWS is empty"
else
    if [ "$SHOULD_START_WORKFLOWS" = "1" ]; then
        ./starter_main && ./worker_main
    else
        ./starter_main
    fi  
fi