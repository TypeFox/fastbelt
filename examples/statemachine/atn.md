# Runtime ATN for statemachine

## Statemachine

```mermaid
flowchart TD
    q0(["SN:0<br/>RuleStart"])
    q1(["SN:1<br/>RuleStop"])
    q10["SN:10<br/>Basic<br/>"]
    q11["SN:12<br/>Basic<br/>"]
    q12["SN:14<br/>Basic<br/>"]
    q13["SN:16<br/>Basic<br/>"]
    q14["SN:17<br/>Basic<br/>"]
    q15{"SN:18<br/>PlusLoopBack<br/><br/>dec=0"}
    q16["SN:19<br/>LoopEnd<br/>"]
    q17["SN:19<br/>Basic<br/>"]
    q18["SN:21<br/>Basic<br/>"]
    q19["SN:22<br/>Basic<br/>"]
    q20{"SN:23<br/>PlusLoopBack<br/><br/>dec=1"}
    q21["SN:24<br/>LoopEnd<br/>"]
    q22["SN:24<br/>Basic<br/>"]
    q23["SN:26<br/>Basic<br/>"]
    q24["SN:28<br/>Basic<br/>"]
    q25["SN:29<br/>Basic<br/>"]
    q26{"SN:30<br/>StarLoopEntry<br/><br/>dec=2"}
    q27["SN:31<br/>LoopEnd<br/>"]
    q28["SN:32<br/>StarLoopBack<br/>"]

    q0 --> q10
    q10 -->|"tok(&quot;statemachine&quot;)"| q11
    q11 -->|"tok(ID)"| q12
    q12 -->|"tok(&quot;events&quot;)"| q13
    q12 --> q16
    q13 -.->|"[Event]"| q14
    q14 --> q15
    q15 --> q13
    q15 --> q16
    q16 --> q17
    q17 -->|"tok(&quot;commands&quot;)"| q18
    q17 --> q21
    q18 -.->|"[Command]"| q19
    q19 --> q20
    q20 --> q18
    q20 --> q21
    q21 --> q22
    q22 -->|"tok(&quot;initialState&quot;)"| q23
    q23 -->|"tok(ID)"| q26
    q24 -.->|"[State]"| q25
    q25 --> q28
    q26 --> q24
    q26 --> q27
    q27 --> q1
    q28 --> q26
```

## Event

```mermaid
flowchart TD
    q2(["SN:2<br/>RuleStart"])
    q3(["SN:3<br/>RuleStop"])
    q29["SN:29<br/>Basic<br/>"]
    q30["SN:30<br/>Basic<br/>"]

    q2 --> q29
    q29 -->|"tok(ID)"| q30
    q30 --> q3
```

## Command

```mermaid
flowchart TD
    q4(["SN:4<br/>RuleStart"])
    q5(["SN:5<br/>RuleStop"])
    q31["SN:31<br/>Basic<br/>"]
    q32["SN:32<br/>Basic<br/>"]

    q4 --> q31
    q31 -->|"tok(ID)"| q32
    q32 --> q5
```

## State

```mermaid
flowchart TD
    q6(["SN:6<br/>RuleStart"])
    q7(["SN:7<br/>RuleStop"])
    q33["SN:33<br/>Basic<br/>"]
    q34["SN:35<br/>Basic<br/>"]
    q35["SN:37<br/>Basic<br/>"]
    q36["SN:39<br/>Basic<br/>"]
    q37["SN:41<br/>Basic<br/>"]
    q38["SN:43<br/>Basic<br/>"]
    q39["SN:44<br/>Basic<br/>"]
    q40["SN:42<br/>Basic<br/>"]
    q41["SN:43<br/>Basic<br/>"]
    q42{"SN:44<br/>StarLoopEntry<br/><br/>dec=3"}
    q43["SN:45<br/>LoopEnd<br/>"]
    q44["SN:46<br/>StarLoopBack<br/>"]
    q45["SN:47<br/>Basic<br/>"]
    q46["SN:48<br/>Basic<br/>"]

    q6 --> q33
    q33 -->|"tok(&quot;state&quot;)"| q34
    q34 -->|"tok(ID)"| q35
    q35 -->|"tok(&quot;actions&quot;)"| q36
    q35 --> q39
    q36 -->|"tok(&quot;{&quot;)"| q37
    q37 -->|"tok(ID)"| q38
    q38 -->|"tok(&quot;}&quot;)"| q39
    q39 --> q42
    q40 -.->|"[Transition]"| q41
    q41 --> q44
    q42 --> q40
    q42 --> q43
    q43 --> q45
    q44 --> q42
    q45 -->|"tok(&quot;end&quot;)"| q46
    q46 --> q7
```

## Transition

```mermaid
flowchart TD
    q8(["SN:8<br/>RuleStart"])
    q9(["SN:9<br/>RuleStop"])
    q47["SN:47<br/>Basic<br/>"]
    q48["SN:49<br/>Basic<br/>"]
    q49["SN:51<br/>Basic<br/>"]
    q50["SN:52<br/>Basic<br/>"]

    q8 --> q47
    q47 -->|"tok(ID)"| q48
    q48 -->|"tok(&quot;=>&quot;)"| q49
    q49 -->|"tok(ID)"| q50
    q50 --> q9
```

