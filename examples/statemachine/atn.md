# Runtime ATN for statemachine

## Statemachine

```mermaid
flowchart TD
    q0(["StateNumber__Statemachine__Start (0)<br/>RuleStart"])
    q1(["StateNumber__Statemachine__Stop (1)<br/>RuleStop"])
    q10["StateNumber__Statemachine_STATEMACHINE (10)<br/>Basic<br/>"]
    q11["StateNumber__Statemachine_Name_ID (11)<br/>Basic<br/>"]
    q12["StateNumber__Statemachine_EVENTS (12)<br/>Basic<br/>"]
    q13["StateNumber__Statemachine__Basic_0 (13)<br/>Basic<br/>"]
    q14["StateNumber__Statemachine__Basic_1 (14)<br/>Basic<br/>"]
    q15{"StateNumber__Statemachine__LoopBack_0 (15)<br/>LoopBack<br/><br/>dec=0"}
    q16["StateNumber__Statemachine__LoopEnd_0 (16)<br/>LoopEnd<br/>"]
    q17{"StateNumber__Statemachine__Basic_2 (17)<br/>Basic<br/><br/>dec=1"}
    q18["StateNumber__Statemachine_COMMANDS (18)<br/>Basic<br/>"]
    q19["StateNumber__Statemachine__Basic_3 (19)<br/>Basic<br/>"]
    q20["StateNumber__Statemachine__Basic_4 (20)<br/>Basic<br/>"]
    q21{"StateNumber__Statemachine__LoopBack_1 (21)<br/>LoopBack<br/><br/>dec=2"}
    q22["StateNumber__Statemachine__LoopEnd_1 (22)<br/>LoopEnd<br/>"]
    q23{"StateNumber__Statemachine__Basic_5 (23)<br/>Basic<br/><br/>dec=3"}
    q24["StateNumber__Statemachine_INITIALSTATE (24)<br/>Basic<br/>"]
    q25["StateNumber__Statemachine_Init_ID (25)<br/>Basic<br/>"]
    q26["StateNumber__Statemachine__Basic_6 (26)<br/>Basic<br/>"]
    q27["StateNumber__Statemachine__Basic_7 (27)<br/>Basic<br/>"]
    q28{"StateNumber__Statemachine__LoopEntry (28)<br/>LoopEntry<br/><br/>dec=4"}
    q29["StateNumber__Statemachine__LoopEnd_2 (29)<br/>LoopEnd<br/>"]
    q30["StateNumber__Statemachine__LoopBack_2 (30)<br/>LoopBack<br/>"]

    q0 --> q10
    q10 -->|"tok(Token_STATEMACHINE)"| q11
    q11 -->|"tok(Token_ID)"| q17
    q12 -->|"tok(Token_EVENTS)"| q13
    q13 -.->|"[Event]"| q14
    q14 --> q15
    q15 --> q13
    q15 --> q16
    q16 --> q23
    q17 --> q12
    q17 --> q16
    q18 -->|"tok(Token_COMMANDS)"| q19
    q19 -.->|"[Command]"| q20
    q20 --> q21
    q21 --> q19
    q21 --> q22
    q22 --> q24
    q23 --> q18
    q23 --> q22
    q24 -->|"tok(Token_INITIALSTATE)"| q25
    q25 -->|"tok(Token_ID)"| q28
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
    q2(["StateNumber__Event__Start (2)<br/>RuleStart"])
    q3(["StateNumber__Event__Stop (3)<br/>RuleStop"])
    q31["StateNumber__Event_Name_ID (31)<br/>Basic<br/>"]
    q32["StateNumber__Event__Basic (32)<br/>Basic<br/>"]

    q2 --> q31
    q31 -->|"tok(Token_ID)"| q32
    q32 --> q3
```

## Command

```mermaid
flowchart TD
    q4(["StateNumber__Command__Start (4)<br/>RuleStart"])
    q5(["StateNumber__Command__Stop (5)<br/>RuleStop"])
    q33["StateNumber__Command_Name_ID (33)<br/>Basic<br/>"]
    q34["StateNumber__Command__Basic (34)<br/>Basic<br/>"]

    q4 --> q33
    q33 -->|"tok(Token_ID)"| q34
    q34 --> q5
```

## State

```mermaid
flowchart TD
    q6(["StateNumber__State__Start (6)<br/>RuleStart"])
    q7(["StateNumber__State__Stop (7)<br/>RuleStop"])
    q35["StateNumber__State_STATE (35)<br/>Basic<br/>"]
    q36["StateNumber__State_Name_ID (36)<br/>Basic<br/>"]
    q37["StateNumber__State_ACTIONS (37)<br/>Basic<br/>"]
    q38["StateNumber__State_LBRACE (38)<br/>Basic<br/>"]
    q39["StateNumber__State_Actions_ID (39)<br/>Basic<br/>"]
    q40["StateNumber__State__Basic_0 (40)<br/>Basic<br/>"]
    q41{"StateNumber__State__LoopBack_0 (41)<br/>LoopBack<br/><br/>dec=5"}
    q42["StateNumber__State__LoopEnd_0 (42)<br/>LoopEnd<br/>"]
    q43["StateNumber__State_RBRACE (43)<br/>Basic<br/>"]
    q44["StateNumber__State__Basic_1 (44)<br/>Basic<br/>"]
    q45{"StateNumber__State__Basic_2 (45)<br/>Basic<br/><br/>dec=6"}
    q46["StateNumber__State__Basic_3 (46)<br/>Basic<br/>"]
    q47["StateNumber__State__Basic_4 (47)<br/>Basic<br/>"]
    q48{"StateNumber__State__LoopEntry (48)<br/>LoopEntry<br/><br/>dec=7"}
    q49["StateNumber__State__LoopEnd_1 (49)<br/>LoopEnd<br/>"]
    q50["StateNumber__State__LoopBack_1 (50)<br/>LoopBack<br/>"]
    q51["StateNumber__State_END (51)<br/>Basic<br/>"]
    q52["StateNumber__State__Basic_5 (52)<br/>Basic<br/>"]

    q6 --> q35
    q35 -->|"tok(Token_STATE)"| q36
    q36 -->|"tok(Token_ID)"| q45
    q37 -->|"tok(Token_ACTIONS)"| q38
    q38 -->|"tok(Token_LBRACE)"| q39
    q39 -->|"tok(Token_ID)"| q40
    q40 --> q41
    q41 --> q39
    q41 --> q42
    q42 --> q43
    q43 -->|"tok(Token_RBRACE)"| q44
    q44 --> q48
    q45 --> q37
    q45 --> q44
    q46 -.->|"[Transition]"| q47
    q47 --> q50
    q48 --> q46
    q48 --> q49
    q49 --> q51
    q50 --> q48
    q51 -->|"tok(Token_END)"| q52
    q52 --> q7
```

## Transition

```mermaid
flowchart TD
    q8(["StateNumber__Transition__Start (8)<br/>RuleStart"])
    q9(["StateNumber__Transition__Stop (9)<br/>RuleStop"])
    q53["StateNumber__Transition_Event_ID (53)<br/>Basic<br/>"]
    q54["StateNumber__Transition_ARROW (54)<br/>Basic<br/>"]
    q55["StateNumber__Transition_State_ID (55)<br/>Basic<br/>"]
    q56["StateNumber__Transition__Basic (56)<br/>Basic<br/>"]

    q8 --> q53
    q53 -->|"tok(Token_ID)"| q54
    q54 -->|"tok(Token_ARROW)"| q55
    q55 -->|"tok(Token_ID)"| q56
    q56 --> q9
```

