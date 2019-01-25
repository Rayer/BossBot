# Keyword based Chatbot

## Abstraction

The engine reacts with user by providing / handling keywords. Engine will highlight keywords, with each keyword it will goes to different state and another sentence or action. If there is no keyword in this scenario, it will be a default action. For example, we can have a dialog like this :

```
Bot : Greetings Rayer, Are you going to [submit] weekly report or you want to [maintain] system?
R : I want to submit report.
Bot : You are not yet submit report this week. What you have [done] or [in development] or [end] of report?
R : Done - MCDS-2219, balabalaba
Bot : Anything else?
R : In development : MCDS-1111
Bot : Anything else?
R : End of report.
Bot : Your report is saved.
```

In this scenario, we can see keywords and actions between bot and user, and we are going to design an engine for interacting between user and chat bot.

## Terms

KBC is basically operated by several scenario and state inside this. KBC will parse keyword and decide it should go to another state of this scenario, or just invoke another scenario.

### Scenario

A scenario means "a statable context" for user. For example, example in abstraction section shows there is several scenario : 

- Main Scenario, the root of all scenarios
- Weekly Report Scenario

It can contains several `state` in a scenario. However, every linar scenario will need to lead directly or indirectly to root scenario -- or by timeout.

A scenario will contains several state, all state shares same context. 

#### Scenario Lifecycle

A scenario lifecycle depends if it is unlinked, and root scenario always lives if user session is active. 

Usually these lifecycle can be managed by 

- `EnterScenario()`
- `ExitScenario()`
- `DisposeScenario()`
- `InvokeWithStrategy()`


#### Scenario Chains

A scenario can call another scenario, and chain(or unhook) together. When a scenario call another scenario, there are several strategy to control whether the lifecycle of callee(and caller) :

- Stack
- Trim
- Replace

If caller is root scenario, only Stack is available when root scenario attempt to create another scenario.

##### Stack

Stack will append next scenario in this chain. For example, when scenario A call scenario B then scenario B call scenario C, it will have this chain :

Root -> A -> B -> C

When scenario C ends, it will re-enter scenario B (and call `EnterScenario()` with argument `reenter : true` and `fromScenario : C`), and so do when scenario B ends.

##### Trim

When original scenario chain is :

Root -> A -> B

and we call scenario C from B with strategy trim, it will become :

Root -> C

In this case, `DisposeScenario()` of scenario A and B will be invoked.

##### Replace

When original scenario chain is :

Root -> A -> B

and we call scenario C from B with strategy trim, it will become :

Root -> A -> C

C will replace B, B's `DisposeScenario()` will be invoked.

### State

It is base operation unit in scenario. Each state can do something like data persistent, calling other state or jump to another scenario. Each scenario have an "entry point" state, and it works like root scenario in another case. 





