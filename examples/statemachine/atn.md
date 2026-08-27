# Runtime ATN for statemachine

## Statemachine

```mermaid
flowchart TD
    q0(["Statemachine__Start (0)<br/>RuleStart"])
    q1(["Statemachine__Stop (1)<br/>RuleStop"])
    q10["Statemachine_statemachine (10)<br/>Basic<br/>"]
    q11["Statemachine_Name_ID (11)<br/>Basic<br/>"]
    q12["Statemachine_events (12)<br/>Basic<br/>"]
    q13["Statemachine__Basic_0 (13)<br/>Basic<br/>"]
    q14["Statemachine__Basic_1 (14)<br/>Basic<br/>"]
    q15{"Statemachine__LoopBack_0 (15)<br/>LoopBack<br/><br/>dec=0"}
    q16["Statemachine__LoopEnd_0 (16)<br/>LoopEnd<br/>"]
    q17{"Statemachine__Basic_2 (17)<br/>Basic<br/><br/>dec=1"}
    q18["Statemachine_commands (18)<br/>Basic<br/>"]
    q19["Statemachine__Basic_3 (19)<br/>Basic<br/>"]
    q20["Statemachine__Basic_4 (20)<br/>Basic<br/>"]
    q21{"Statemachine__LoopBack_1 (21)<br/>LoopBack<br/><br/>dec=2"}
    q22["Statemachine__LoopEnd_1 (22)<br/>LoopEnd<br/>"]
    q23{"Statemachine__Basic_5 (23)<br/>Basic<br/><br/>dec=3"}
    q24["Statemachine_initialState (24)<br/>Basic<br/>"]
    q25["Statemachine_Init_ID (25)<br/>Basic<br/>"]
    q26["Statemachine__Basic_6 (26)<br/>Basic<br/>"]
    q27["Statemachine__Basic_7 (27)<br/>Basic<br/>"]
    q28{"Statemachine__LoopEntry (28)<br/>LoopEntry<br/><br/>dec=4"}
    q29["Statemachine__LoopEnd_2 (29)<br/>LoopEnd<br/>"]
    q30["Statemachine__LoopBack_2 (30)<br/>LoopBack<br/>"]

    q0 --> q10
    q10 -->|"tok(&quot;statemachine&quot;)"| q11
    q11 -->|"tok(ID)"| q17
    q12 -->|"tok(&quot;events&quot;)"| q13
    q13 -.->|"[Event]"| q14
    q14 --> q15
    q15 --> q13
    q15 --> q16
    q16 --> q23
    q17 --> q12
    q17 --> q16
    q18 -->|"tok(&quot;commands&quot;)"| q19
    q19 -.->|"[Command]"| q20
    q20 --> q21
    q21 --> q19
    q21 --> q22
    q22 --> q24
    q23 --> q18
    q23 --> q22
    q24 -->|"tok(&quot;initialState&quot;)"| q25
    q25 -->|"tok(ID)"| q28
    q26 -.->|"[State]"| q27
    q27 --> q30
    q28 --> q26
    q28 --> q29
    q29 --> q1
    q30 --> q28
```

## Event

```mermaid
flowchart TD
    q2(["Event__Start (2)<br/>RuleStart"])
    q3(["Event__Stop (3)<br/>RuleStop"])
    q31["Event_Name_ID (31)<br/>Basic<br/>"]
    q32["Event__Basic (32)<br/>Basic<br/>"]

    q2 --> q31
    q31 -->|"tok(ID)"| q32
    q32 --> q3
```

## Command

```mermaid
flowchart TD
    q4(["Command__Start (4)<br/>RuleStart"])
    q5(["Command__Stop (5)<br/>RuleStop"])
    q33["Command_Name_ID (33)<br/>Basic<br/>"]
    q34["Command__Basic (34)<br/>Basic<br/>"]

    q4 --> q33
    q33 -->|"tok(ID)"| q34
    q34 --> q5
```

## State

```mermaid
flowchart TD
    q6(["State__Start (6)<br/>RuleStart"])
    q7(["State__Stop (7)<br/>RuleStop"])
    q35["State_state (35)<br/>Basic<br/>"]
    q36["State_Name_ID (36)<br/>Basic<br/>"]
    q37["State_actions (37)<br/>Basic<br/>"]
    q38["State_LeftBrace (38)<br/>Basic<br/>"]
    q39["State_Actions_ID (39)<br/>Basic<br/>"]
    q40["State__Basic_0 (40)<br/>Basic<br/>"]
    q41{"State__LoopBack_0 (41)<br/>LoopBack<br/><br/>dec=5"}
    q42["State__LoopEnd_0 (42)<br/>LoopEnd<br/>"]
    q43["State_RightBrace (43)<br/>Basic<br/>"]
    q44["State__Basic_1 (44)<br/>Basic<br/>"]
    q45{"State__Basic_2 (45)<br/>Basic<br/><br/>dec=6"}
    q46["State__Basic_3 (46)<br/>Basic<br/>"]
    q47["State__Basic_4 (47)<br/>Basic<br/>"]
    q48{"State__LoopEntry (48)<br/>LoopEntry<br/><br/>dec=7"}
    q49["State__LoopEnd_1 (49)<br/>LoopEnd<br/>"]
    q50["State__LoopBack_1 (50)<br/>LoopBack<br/>"]
    q51["State_end (51)<br/>Basic<br/>"]
    q52["State__Basic_5 (52)<br/>Basic<br/>"]

    q6 --> q35
    q35 -->|"tok(&quot;state&quot;)"| q36
    q36 -->|"tok(ID)"| q45
    q37 -->|"tok(&quot;actions&quot;)"| q38
    q38 -->|"tok(&quot;{&quot;)"| q39
    q39 -->|"tok(ID)"| q40
    q40 --> q41
    q41 --> q39
    q41 --> q42
    q42 --> q43
    q43 -->|"tok(&quot;}&quot;)"| q44
    q44 --> q48
    q45 --> q37
    q45 --> q44
    q46 -.->|"[Transition]"| q47
    q47 --> q50
    q48 --> q46
    q48 --> q49
    q49 --> q51
    q50 --> q48
    q51 -->|"tok(&quot;end&quot;)"| q52
    q52 --> q7
```

## Transition

```mermaid
flowchart TD
    q8(["Transition__Start (8)<br/>RuleStart"])
    q9(["Transition__Stop (9)<br/>RuleStop"])
    q53["Transition_Event_ID (53)<br/>Basic<br/>"]
    q54["Transition_EqualsGreaterThan (54)<br/>Basic<br/>"]
    q55["Transition_State_ID (55)<br/>Basic<br/>"]
    q56["Transition__Basic (56)<br/>Basic<br/>"]

    q8 --> q53
    q53 -->|"tok(ID)"| q54
    q54 -->|"tok(&quot;=>&quot;)"| q55
    q55 -->|"tok(ID)"| q56
    q56 --> q9
```

