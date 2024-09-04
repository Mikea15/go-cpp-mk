// Copyright (c) 2022-2024 Michael Adaixo

#pragma once

#include "CoreMinimal.h"
#include "FlowActorReference.h"
#include "FlowTypes.h"

#include "Templates/Function.h"

#if WITH_EDITOR
#include "Misc/DataValidation.h"
#endif

#include "FlowPilotTask.generated.h"

#if !UE_BUILD_SHIPPING && !UE_BUILD_TEST
struct FDebugUtils;
#endif // !UE_BUILD_SHIPPING && !UE_BUILD_TEST

struct FFlowContext;
struct FFlowActorReference;
class UFlowPilotComponent;
class UFlowPilotParent;

#if WITH_EDITOR
struct FFlowPilotTaskEditorData
{
	bool bExpanded = false;
};
#endif

/**
 * FlowPilotTask
 * Base class for any task that can be run by FlowPilotComponent
 * Class is tickable.
 * If Tick not implemented, will automatically succeed on first tick
 */
UCLASS(Abstract, EditInlineNew, DefaultToInstanced, CollapseCategories, AutoExpandCategories=("FlowPilot"))
class FLOWPILOT_API UFlowPilotTask : public UObject
{
	GENERATED_BODY()
	
public:
	UFlowPilotTask();
	
	// UFlowPilotTask
	// Setups Tasks. Called once per FlowPilotExecution, even after restarts.
	virtual void Setup(FFlowContext* InContext);
	
	// Called when starting this Task. Returns true on success
	virtual bool Enter();
	
	// Called on Tick. Will success automatically if not implemented by Child classes
	virtual EFPTaskResult Tick(float DeltaTime);
	
	// Called when Task Finished
	virtual void Exit(EFPTaskResult TaskResult);
	
	// Resets all Tasks into their Setup States
	virtual void Reset();
	//~UFlowPilotTask

	// Disabled Tasks are skipped during execution
	bool IsEnabled() const { return bEnabled; }
	
	// Enables or Disables Task. Disabled Tasks will be skipped.
	void SetIsEnabled(bool bInEnabled) { bEnabled = bInEnabled; }

	// Returns Task Name
	FName GetTaskName() const { return TaskName; }

	// Sets Task Name
	void SetTaskName(FName NewTaskName) { TaskName = NewTaskName; }

	// Get Task Description
	FName GetTaskDescription() const { return Description; }

	// Returns True if Task has Parent Task.
	// Returns False if Task is Root Sequence Task
	bool HasParent() const;
	
	// Returns Parent Task or nullptr
	UFlowPilotTask* GetParent() const;

	// Sets Parent Task
	void SetParent(UFlowPilotTask* InParent) { Parent = InParent; }

	// Returns True if This task is a FlowPilotParent Task containing children Tasks.
	bool IsParent() const;
	
	// Returns this Cast to FlowPilotParent task.
	UFlowPilotParent* GetAsParent();
	

#if WITH_EDITOR
	// Returns true if valid. Child Tasks should implement their Validations
	virtual bool IsTaskDataValid(FDataValidationContext& InContext) { return true; }

	virtual FName GetBrush() const;
#endif

#if !UE_BUILD_SHIPPING && !UE_BUILD_TEST
	// Gathers information to display to debug view about Task.
	virtual void GetRuntimeDescription(TArray<FString>& OutLines) const {};
#endif
	//~UFlowPilotTask
	
	// Returns true when Task Started
	UFUNCTION(BlueprintCallable, Category="FlowPilotTask")
	bool HasStarted() const;

	// Returns true when Task in Progress and Not Complete
	UFUNCTION(BlueprintCallable, Category="FlowPilotTask")
	bool IsActive() const;

	// Returns true when Task is Complete
	UFUNCTION(BlueprintCallable, Category="FlowPilotTask")
	bool IsComplete() const;

	// Executes 'InFunc' to all Actors found from 'ActorReference'
	bool ForEachActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(AActor* const /*Actor*/)> InFunc) const;

	// Executes 'InFunc' to all const Actors found from 'ActorReference'
	// Const means the function should not modify 'Actors'
	bool ForEachConstActor(const FFlowActorReference& ActorReference, TFunctionRef<bool(const AActor* const /*Actor*/)> InFunc) const;

#if WITH_EDITOR
	FFlowPilotTaskEditorData EditorData;
#endif

protected:
	// Returns FlowPilotComponent
	UFUNCTION(BlueprintCallable, Category = "FlowPilotTask")
	UFlowPilotComponent* GetFlowPilotComponent() const;

	// Returns FlowPilotComponent Owner Actor
	UFUNCTION(BlueprintCallable, Category = "FlowPilotTask")
	AActor* GetFlowPilotOwnerActor() const;

	// Returns FlowPilot Actor World
	UFUNCTION(BlueprintCallable, Category = "FlowPilotTask")
	UWorld* GetWorldContext() const;

	virtual UWorld* GetWorld() const override;

protected:
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

	const FFlowContext* Context = nullptr;
	EFPInternalTaskState InternalState;

#if WITH_EDITOR
	friend class SFlowPilotViewRow;
	friend class FFlowPilotViewModel; 
#endif
};

// Class 2
UCLASS(Abstract, EditInlineNew, DefaultToInstanced, CollapseCategories, AutoExpandCategories=("FlowPilot"))
class FLOWPILOT_API UFlowPilotTask2 : public UObject, USomeOtherClass, Interface
{
	// Setups Tasks. Called once per FlowPilotExecution, even after restarts.
	virtual void Setup(FFlowContext* InContext);
	
	// Called when starting this Task. Returns true on success
	virtual bool Enter();
};

// Class 3
UCLASS(Abstract, EditInlineNew, DefaultToInstanced, CollapseCategories, AutoExpandCategories=("FlowPilot"))
class FLOWPILOT_API UFlowPilotTask3 : public UObject
{
	// Setups Tasks. Called once per FlowPilotExecution, even after restarts.
	virtual void Setup(FFlowContext* InContext);
	
	// Called when starting this Task. Returns true on success
	virtual bool Enter();
};
