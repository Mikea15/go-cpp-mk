---
title: example.h
description: Reference page for example.h
---

import Ref from '../../../components/Ref.astro'

import { Image } from 'astro:assets';
import schema from '../../../assets/uflowpilottask_executionflow.png';

<Ref c="UFlowPilotTask" p="UObject" />

## Description

<Image src={schema} alt="UFlowPilotTask Schema" />

`UFlowPilotTask` is the base class for all Tasks that can run in FlowPilot.

The schema above shows succintely how it works.

Each Tasks starts its execution by calling their `Enter()` methods. They can either succeed or fail at this stage.
Upon failing, the task ends. If succeeded, `Tick` is then called and we check the Task Result ( InProgress, Succeeded, Failed, Error ).
If the Task returns `Succeeded` the `Exit()` method is called.

Tasks can also be cancelled, and end directly.

Each class that can be implemented need to override a couple of virtual methods exposed by this task, either in Cpp or implementing the Blueprint versions when creating a `UFlowPilotTask` via blueprint.

## File Info

__FileName:__ `example.h`
- __Class List:__ 
[ [`UFlowPilotTask`](#UFlowPilotTask) | [`UFlowPilotTask2`](#UFlowPilotTask2) | [`UFlowPilotTask3`](#UFlowPilotTask3) ]
- __Struct List:__ 
[ [`FFlowActorReference`](#FFlowActorReference) ]
- __Enum List:__ 
[ [`EFPInternalTaskState`](#EFPInternalTaskState) | [`EFPStopType`](#EFPStopType) | [`EFPTaskResult`](#EFPTaskResult) ]


## `EFPInternalTaskState` 


### Properties

```cpp
// Not started yet 
Invalid = 0,

// Setup done 
Setup,

// Started Execution 
Started,

// Is in Progress 
Ticking,

// Completed 
Completed

```


## `EFPStopType` 


### Properties

```cpp
// Cancel Execution will call Exit with Failure state on Running Tasks 
CancelExecution,

// Stop Now will just stop execution. Exit won't be called. 
StopNow

```


## `EFPTaskResult` 


### Properties

```cpp
// Not started yet 
None,

// In progress and Ticking 
InProgress,

// Complete with Success Result 
Succeeded,

// Complete with Fail Result 
Failed,

// Not Complete. Return Error 
Error

```


## `FFlowActorReference` 


 \
FFlowActorReference 

### Properties

```cpp
// What's the scope of the Subject Actor 
// (Self, In Level or Runtime) 
UPROPERTY(EditAnywhere, Category = "Actor Reference")
EFlowActorScope Scope = EFlowActorScope::Self;

// Level Actor Reference 
UPROPERTY(EditAnywhere, Category = "Actor Reference", meta=(EditCondition="Scope==EFlowActorScope::InLevel", EditConditionHides))
TSoftObjectPtr<AActor> LevelActor;

// Runtime Actor, found via Gameplay Tag. 
UPROPERTY(EditAnywhere, Category = "Actor Reference", meta=(EditCondition="Scope==EFlowActorScope::Runtime", EditConditionHides))
FGameplayTag ExternalActorTag;

```


## `UFlowPilotTask` 


__Parent Classes:__
[ `UObject` ]

 \
FlowPilotTask \
Base class for any task that can be run by FlowPilotComponent \
Class is tickable. \
If Tick not implemented, will automatically succeed on first tick 

### Properties

```cpp
// Task name 
UPROPERTY(EditDefaultsOnly, Category = "Task Options")
FName TaskName = {};

// Task description 
UPROPERTY(EditDefaultsOnly, Category = "Task Options")
FName Description = {};

// If False, Will skip this Task's execution 
UPROPERTY(EditDefaultsOnly, Category = "Task Options")
uint8 bEnabled: 1;

// Parent Task 
UPROPERTY()
UFlowPilotTask* Parent = nullptr;

```

### Functions

#### `Setup`
> Setups Tasks. Called once per FlowPilotExecution, even after restarts. 
```cpp
virtual void Setup(FFlowContext* InContext);
```
#### `Enter`
> Called when starting this Task. Returns true on success 
```cpp
virtual bool Enter();
```
#### `Tick`
> Called on Tick. Will success automatically if not implemented by Child classes 
```cpp
virtual EFPTaskResult Tick(float DeltaTime);
```
#### `Exit`
> Called when Task Finished 
```cpp
virtual void Exit(EFPTaskResult TaskResult);
```
#### `Reset`
> Resets all Tasks into their Setup States 
```cpp
virtual void Reset();
```
#### `IsEnabled`
> Disabled Tasks are skipped during execution 
```cpp
bool IsEnabled() const { return bEnabled; }
```
#### `SetIsEnabled`
> Enables or Disables Task. Disabled Tasks will be skipped. 
```cpp
void SetIsEnabled(bool bInEnabled) { bEnabled = bInEnabled; }
```
#### `GetTaskName`
> Returns Task Name 
```cpp
FName GetTaskName() const { return TaskName; }
```
#### `SetTaskName`
> Sets Task Name 
```cpp
void SetTaskName(FName NewTaskName) { TaskName = NewTaskName; }
```
#### `GetTaskDescription`
> Get Task Description 
```cpp
FName GetTaskDescription() const { return Description; }
```
#### `HasParent`
> Returns True if Task has Parent Task. \
> Returns False if Task is Root Sequence Task 
```cpp
bool HasParent() const;
```
#### `GetParent`
> Returns Parent Task or nullptr 
```cpp
UFlowPilotTask* GetParent() const;
```
#### `SetParent`
> Sets Parent Task 
```cpp
void SetParent(UFlowPilotTask* InParent) { Parent = InParent; }
```
#### `IsParent`
> Returns True if This task is a FlowPilotParent Task containing children Tasks. 
```cpp
bool IsParent() const;
```
#### `GetAsParent`
> Returns this Cast to FlowPilotParent task. 
```cpp
UFlowPilotParent* GetAsParent();
```
#### `HasStarted`
> Returns true when Task Started 
```cpp
bool HasStarted() const;
```
#### `IsActive`
> Returns true when Task in Progress and Not Complete 
```cpp
bool IsActive() const;
```
#### `IsComplete`
> Returns true when Task is Complete 
```cpp
bool IsComplete() const;
```
#### `ForEachActor`
> Executes 'InFunc' to all Actors found from 'ActorReference' 
```cpp
bool ForEachActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(AActor* const /*Actor*/)> InFunc) const;
```
#### `ForEachConstActor`
> Executes 'InFunc' to all const Actors found from 'ActorReference' \
> Const means the function should not modify 'Actors' 
```cpp
bool ForEachConstActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(const AActor* const /*Actor*/)> InFunc) const;
```
#### `GetFlowPilotComponent`
> Returns FlowPilotComponent 
```cpp
UFlowPilotComponent* GetFlowPilotComponent() const;
```
#### `GetFlowPilotOwnerActor`
> Returns FlowPilotComponent Owner Actor 
```cpp
AActor* GetFlowPilotOwnerActor() const;
```
#### `GetWorldContext`
> Returns FlowPilot Actor World 
```cpp
UWorld* GetWorldContext() const;
```


## `UFlowPilotTask2` 


__Parent Classes:__
[ `UObject,`, `USomeOtherClass,`, `Interface` ]

Class 2 

### Functions

#### `Setup`
> Setups Tasks. Called once per FlowPilotExecution, even after restarts. 
```cpp
virtual void Setup(FFlowContext* InContext);
```
#### `Enter`
> Called when starting this Task. Returns true on success 
```cpp
virtual bool Enter();
```


## `UFlowPilotTask3` 


__Parent Classes:__
[ `UObject` ]

Class 3 

### Functions

#### `Setup`
> Setups Tasks. Called once per FlowPilotExecution, even after restarts. 
```cpp
virtual void Setup(FFlowContext* InContext);
```
#### `Enter`
> Called when starting this Task. Returns true on success 
```cpp
virtual bool Enter();
```
