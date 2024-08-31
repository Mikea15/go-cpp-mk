# `UFlowPilotTask` : `UObject`

## Properties

| Property | Description |
|----------|-------------|
| `UPROPERTY(EditDefaultsOnly, Category = "Task Options")`<br>`FName TaskName = {};` |  |
| `UPROPERTY(EditDefaultsOnly, Category = "Task Options")`<br>`FName Description = {};` |  |
| `UPROPERTY(EditDefaultsOnly, Category = "Task Options")`<br>`uint8 bEnabled: 1;` | If False, Will skip this Task's execution |
| `UPROPERTY()`<br>`UFlowPilotTask* Parent = nullptr;` |  |

## Methods

```cpp
UFlowPilotTask();
```

/**
* FlowPilotTask
* Base class for any task that can be run by FlowPilotComponent
* Class is tickable.
* If Tick not implemented, will automatically succeed on first tick
*/

```cpp
virtual void Setup(FFlowContext* InContext);
```

 Setups Tasks. Called once per FlowPilotExecution, even after restarts.

```cpp
virtual bool Enter();
```

 Called when starting this Task. Returns true on success

```cpp
virtual EFPTaskResult Tick(float DeltaTime);
```

 Called on Tick. Will success automatically if not implemented by Child classes

```cpp
virtual void Exit(EFPTaskResult TaskResult);
```

 Called when Task Finished

```cpp
virtual void Reset();
```

 Resets all Tasks into their Setup States

```cpp
bool HasParent() const;
```

 Disabled Tasks are skipped during execution
 Enables or Disables Task. Disabled Tasks will be skipped.
 Returns Task Name
 Sets Task Name
 Get Task Description
 Returns True if Task has Parent Task.
 Returns False if Task is Root Sequence Task

```cpp
UFlowPilotTask* GetParent() const;
```

 Returns Parent Task or nullptr

```cpp
bool IsParent() const;
```

 Sets Parent Task
 Returns True if This task is a FlowPilotParent Task containing children Tasks.

```cpp
UFlowPilotParent* GetAsParent();
```

 Returns this Cast to FlowPilotParent task.

```cpp
virtual FName GetBrush() const;
```

 Returns true if valid. Child Tasks should implement their Validations

```cpp
virtual void GetRuntimeDescription(TArray<FString>& OutLines) const {};
```

 Gathers information to display to debug view about Task.

```cpp
bool HasStarted() const;
```

 Returns true when Task Started

```cpp
bool IsActive() const;
```

 Returns true when Task in Progress and Not Complete

```cpp
bool IsComplete() const;
```

 Returns true when Task is Complete

```cpp
bool ForEachActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(AActor* const /*Actor*/)> InFunc) const;
```

 Executes 'InFunc' to all Actors found from 'ActorReference'

```cpp
bool ForEachConstActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(const AActor* const /*Actor*/)> InFunc) const;
```

 Executes 'InFunc' to all const Actors found from 'ActorReference'
 Const means the function should not modify 'Actors'

```cpp
UFlowPilotComponent* GetFlowPilotComponent() const;
```

 Returns FlowPilotComponent

```cpp
AActor* GetFlowPilotOwnerActor() const;
```

 Returns FlowPilotComponent Owner Actor

```cpp
UWorld* GetWorldContext() const;
```

 Returns FlowPilot Actor World

```cpp
virtual UWorld* GetWorld() const override;
```


