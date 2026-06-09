# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["Grammar__Start (0)<br/>RuleStart"])
    q1(["Grammar__Stop (1)<br/>RuleStop"])
    q50["Grammar_grammar (50)<br/>Basic<br/>"]
    q51["Grammar_Name_ID (51)<br/>Basic<br/>"]
    q52["Grammar_Semicolon (52)<br/>Basic<br/>"]
    q53["Grammar__Basic_0 (53)<br/>Basic<br/>"]
    q54["Grammar__Basic_1 (54)<br/>Basic<br/>"]
    q55["Grammar__Basic_2 (55)<br/>Basic<br/>"]
    q56["Grammar__Basic_3 (56)<br/>Basic<br/>"]
    q57["Grammar__Basic_4 (57)<br/>Basic<br/>"]
    q58["Grammar__Basic_5 (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_6 (59)<br/>Basic<br/>"]
    q60["Grammar__Basic_7 (60)<br/>Basic<br/>"]
    q61["Grammar__Basic_8 (61)<br/>Basic<br/>"]
    q62["Grammar__BlockEnd (62)<br/>BlockEnd<br/>"]
    q63{"Grammar__LoopEntry (63)<br/>LoopEntry<br/><br/>dec=0"}
    q64["Grammar__LoopEnd (64)<br/>LoopEnd<br/>"]
    q65["Grammar__LoopBack (65)<br/>LoopBack<br/>"]

    q0 --> q50
    q50 -->|"tok(&quot;grammar&quot;)"| q51
    q51 -->|"tok(ID)"| q52
    q52 -->|"tok(&quot;;&quot;)"| q63
    q53 -.->|"[ParserRule]"| q54
    q54 --> q62
    q55 -.->|"[Token]"| q56
    q56 --> q62
    q57 -.->|"[Interface]"| q58
    q58 --> q62
    q59 -.->|"[CompositeRule]"| q60
    q60 --> q62
    q61 --> q53
    q61 --> q55
    q61 --> q57
    q61 --> q59
    q62 --> q65
    q63 --> q61
    q63 --> q64
    q64 --> q1
    q65 --> q63
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q66["Interface_interface (66)<br/>Basic<br/>"]
    q67["Interface_Name_ID (67)<br/>Basic<br/>"]
    q68["Interface_extends (68)<br/>Basic<br/>"]
    q69["Interface_Extends_ID_0 (69)<br/>Basic<br/>"]
    q70["Interface_Comma (70)<br/>Basic<br/>"]
    q71["Interface_Extends_ID_1 (71)<br/>Basic<br/>"]
    q72["Interface__Basic_0 (72)<br/>Basic<br/>"]
    q73{"Interface__LoopEntry_0 (73)<br/>LoopEntry<br/><br/>dec=1"}
    q74["Interface__LoopEnd_0 (74)<br/>LoopEnd<br/>"]
    q75["Interface__LoopBack_0 (75)<br/>LoopBack<br/>"]
    q76["Interface_LeftBrace (76)<br/>Basic<br/>"]
    q77["Interface__Basic_1 (77)<br/>Basic<br/>"]
    q78["Interface__Basic_2 (78)<br/>Basic<br/>"]
    q79{"Interface__LoopEntry_1 (79)<br/>LoopEntry<br/><br/>dec=2"}
    q80["Interface__LoopEnd_1 (80)<br/>LoopEnd<br/>"]
    q81["Interface__LoopBack_1 (81)<br/>LoopBack<br/>"]
    q82["Interface_RightBrace (82)<br/>Basic<br/>"]
    q83["Interface__Basic_3 (83)<br/>Basic<br/>"]

    q2 --> q66
    q66 -->|"tok(&quot;interface&quot;)"| q67
    q67 -->|"tok(ID)"| q68
    q68 -->|"tok(&quot;extends&quot;)"| q69
    q68 --> q74
    q69 -->|"tok(ID)"| q73
    q70 -->|"tok(&quot;,&quot;)"| q71
    q71 -->|"tok(ID)"| q72
    q72 --> q75
    q73 --> q70
    q73 --> q74
    q74 --> q76
    q75 --> q73
    q76 -->|"tok(&quot;{&quot;)"| q79
    q77 -.->|"[Field]"| q78
    q78 --> q81
    q79 --> q77
    q79 --> q80
    q80 --> q82
    q81 --> q79
    q82 -->|"tok(&quot;}&quot;)"| q83
    q83 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q84["Field_Name_ID (84)<br/>Basic<br/>"]
    q85["Field__Basic_0 (85)<br/>Basic<br/>"]
    q86["Field__Basic_1 (86)<br/>Basic<br/>"]

    q4 --> q84
    q84 -->|"tok(ID)"| q85
    q85 -.->|"[FieldType]"| q86
    q86 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q87["FieldType__Basic_0 (87)<br/>Basic<br/>"]
    q88["FieldType__Basic_1 (88)<br/>Basic<br/>"]
    q89["FieldType__Basic_2 (89)<br/>Basic<br/>"]
    q90["FieldType__Basic_3 (90)<br/>Basic<br/>"]
    q91["FieldType__Basic_4 (91)<br/>Basic<br/>"]
    q92["FieldType__Basic_5 (92)<br/>Basic<br/>"]
    q93["FieldType__Basic_6 (93)<br/>Basic<br/>"]
    q94["FieldType__Basic_7 (94)<br/>Basic<br/>"]
    q95["FieldType__Basic_8 (95)<br/>Basic<br/>"]
    q96["FieldType__BlockEnd (96)<br/>BlockEnd<br/>"]

    q6 --> q95
    q87 -.->|"[SimpleType]"| q88
    q88 --> q96
    q89 -.->|"[ReferenceType]"| q90
    q90 --> q96
    q91 -.->|"[ArrayType]"| q92
    q92 --> q96
    q93 -.->|"[PrimitiveType]"| q94
    q94 --> q96
    q95 --> q87
    q95 --> q89
    q95 --> q91
    q95 --> q93
    q96 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q97["ArrayType_LeftBracket (97)<br/>Basic<br/>"]
    q98["ArrayType_RightBracket (98)<br/>Basic<br/>"]
    q99["ArrayType__Basic_0 (99)<br/>Basic<br/>"]
    q100["ArrayType__Basic_1 (100)<br/>Basic<br/>"]

    q8 --> q97
    q97 -->|"tok(&quot;[&quot;)"| q98
    q98 -->|"tok(&quot;]&quot;)"| q99
    q99 -.->|"[FieldType]"| q100
    q100 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q101["ReferenceType_Asterisk (101)<br/>Basic<br/>"]
    q102["ReferenceType_Type_ID (102)<br/>Basic<br/>"]
    q103["ReferenceType__Basic (103)<br/>Basic<br/>"]

    q10 --> q101
    q101 -->|"tok(&quot;*&quot;)"| q102
    q102 -->|"tok(ID)"| q103
    q103 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q104["SimpleType_Type_ID (104)<br/>Basic<br/>"]
    q105["SimpleType__Basic (105)<br/>Basic<br/>"]

    q12 --> q104
    q104 -->|"tok(ID)"| q105
    q105 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q106["PrimitiveType_Type_string (106)<br/>Basic<br/>"]
    q107["PrimitiveType__Basic_0 (107)<br/>Basic<br/>"]
    q108["PrimitiveType_Type_bool (108)<br/>Basic<br/>"]
    q109["PrimitiveType__Basic_1 (109)<br/>Basic<br/>"]
    q110["PrimitiveType_Type_composite (110)<br/>Basic<br/>"]
    q111["PrimitiveType__Basic_2 (111)<br/>Basic<br/>"]
    q112["PrimitiveType__Basic_3 (112)<br/>Basic<br/>"]
    q113["PrimitiveType__BlockEnd (113)<br/>BlockEnd<br/>"]

    q14 --> q112
    q106 -->|"tok(&quot;string&quot;)"| q107
    q107 --> q113
    q108 -->|"tok(&quot;bool&quot;)"| q109
    q109 --> q113
    q110 -->|"tok(&quot;composite&quot;)"| q111
    q111 --> q113
    q112 --> q106
    q112 --> q108
    q112 --> q110
    q113 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q114["ParserRule_Name_ID (114)<br/>Basic<br/>"]
    q115["ParserRule_returns (115)<br/>Basic<br/>"]
    q116["ParserRule_ReturnType_ID (116)<br/>Basic<br/>"]
    q117["ParserRule__Basic_0 (117)<br/>Basic<br/>"]
    q118["ParserRule_Colon (118)<br/>Basic<br/>"]
    q119["ParserRule__Basic_1 (119)<br/>Basic<br/>"]
    q120["ParserRule_Semicolon (120)<br/>Basic<br/>"]
    q121["ParserRule__Basic_2 (121)<br/>Basic<br/>"]

    q16 --> q114
    q114 -->|"tok(ID)"| q115
    q115 -->|"tok(&quot;returns&quot;)"| q116
    q115 --> q117
    q116 -->|"tok(ID)"| q117
    q117 --> q118
    q118 -->|"tok(&quot;:&quot;)"| q119
    q119 -.->|"[Alternatives]"| q120
    q120 -->|"tok(&quot;;&quot;)"| q121
    q121 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q122["Token_Type_hidden (122)<br/>Basic<br/>"]
    q123["Token__Basic_0 (123)<br/>Basic<br/>"]
    q124["Token_Type_comment (124)<br/>Basic<br/>"]
    q125["Token__Basic_1 (125)<br/>Basic<br/>"]
    q126["Token__Basic_2 (126)<br/>Basic<br/>"]
    q127["Token__BlockEnd (127)<br/>BlockEnd<br/>"]
    q128["Token_token (128)<br/>Basic<br/>"]
    q129["Token_Name_ID (129)<br/>Basic<br/>"]
    q130["Token_Colon (130)<br/>Basic<br/>"]
    q131["Token_Regexp_RegexLiteral (131)<br/>Basic<br/>"]
    q132["Token_Semicolon (132)<br/>Basic<br/>"]
    q133["Token__Basic_3 (133)<br/>Basic<br/>"]

    q18 --> q126
    q122 -->|"tok(&quot;hidden&quot;)"| q123
    q123 --> q127
    q124 -->|"tok(&quot;comment&quot;)"| q125
    q125 --> q127
    q126 --> q122
    q126 --> q124
    q126 --> q127
    q127 --> q128
    q128 -->|"tok(&quot;token&quot;)"| q129
    q129 -->|"tok(ID)"| q130
    q130 -->|"tok(&quot;:&quot;)"| q131
    q131 -->|"tok(RegexLiteral)"| q132
    q132 -->|"tok(&quot;;&quot;)"| q133
    q133 --> q19
```

## Alternatives

```mermaid
flowchart TD
    q20(["Alternatives__Start (20)<br/>RuleStart"])
    q21(["Alternatives__Stop (21)<br/>RuleStop"])
    q134["Alternatives__Basic_0 (134)<br/>Basic<br/>"]
    q135["Alternatives_Pipe (135)<br/>Basic<br/>"]
    q136["Alternatives__Basic_1 (136)<br/>Basic<br/>"]
    q137["Alternatives__Basic_2 (137)<br/>Basic<br/>"]
    q138{"Alternatives__LoopBack (138)<br/>LoopBack<br/><br/>dec=3"}
    q139["Alternatives__LoopEnd (139)<br/>LoopEnd<br/>"]

    q20 --> q134
    q134 -.->|"[Group]"| q135
    q135 -->|"tok(&quot;|&quot;)"| q136
    q135 --> q139
    q136 -.->|"[Group]"| q137
    q137 --> q138
    q138 --> q135
    q138 --> q139
    q139 --> q21
```

## Group

```mermaid
flowchart TD
    q22(["Group__Start (22)<br/>RuleStart"])
    q23(["Group__Stop (23)<br/>RuleStop"])
    q140["Group__Basic_0 (140)<br/>Basic<br/>"]
    q141["Group__Basic_1 (141)<br/>Basic<br/>"]
    q142["Group__Basic_2 (142)<br/>Basic<br/>"]
    q143{"Group__LoopBack (143)<br/>LoopBack<br/><br/>dec=4"}
    q144["Group__LoopEnd (144)<br/>LoopEnd<br/>"]

    q22 --> q140
    q140 -.->|"[Element]"| q141
    q141 -.->|"[Element]"| q142
    q141 --> q144
    q142 --> q143
    q143 --> q141
    q143 --> q144
    q144 --> q23
```

## Element

```mermaid
flowchart TD
    q24(["Element__Start (24)<br/>RuleStart"])
    q25(["Element__Stop (25)<br/>RuleStop"])
    q145["Element__Basic_0 (145)<br/>Basic<br/>"]
    q146["Element__Basic_1 (146)<br/>Basic<br/>"]
    q147["Element__Basic_2 (147)<br/>Basic<br/>"]
    q148["Element__Basic_3 (148)<br/>Basic<br/>"]
    q149["Element__Basic_4 (149)<br/>Basic<br/>"]
    q150["Element__Basic_5 (150)<br/>Basic<br/>"]
    q151["Element__Basic_6 (151)<br/>Basic<br/>"]
    q152["Element__Basic_7 (152)<br/>Basic<br/>"]
    q153["Element_LeftParen (153)<br/>Basic<br/>"]
    q154["Element__Basic_8 (154)<br/>Basic<br/>"]
    q155["Element_RightParen (155)<br/>Basic<br/>"]
    q156["Element__Basic_9 (156)<br/>Basic<br/>"]
    q157["Element__Basic_10 (157)<br/>Basic<br/>"]
    q158["Element__BlockEnd_0 (158)<br/>BlockEnd<br/>"]
    q159["Element_Cardinality_Asterisk (159)<br/>Basic<br/>"]
    q160["Element__Basic_11 (160)<br/>Basic<br/>"]
    q161["Element_Cardinality_Plus (161)<br/>Basic<br/>"]
    q162["Element__Basic_12 (162)<br/>Basic<br/>"]
    q163["Element_Cardinality_Question (163)<br/>Basic<br/>"]
    q164["Element__Basic_13 (164)<br/>Basic<br/>"]
    q165["Element__Basic_14 (165)<br/>Basic<br/>"]
    q166["Element__BlockEnd_1 (166)<br/>BlockEnd<br/>"]

    q24 --> q157
    q145 -.->|"[Keyword]"| q146
    q146 --> q158
    q147 -.->|"[Assignment]"| q148
    q148 --> q158
    q149 -.->|"[RuleCall]"| q150
    q150 --> q158
    q151 -.->|"[Action]"| q152
    q152 --> q158
    q153 -->|"tok(&quot;(&quot;)"| q154
    q154 -.->|"[Alternatives]"| q155
    q155 -->|"tok(&quot;)&quot;)"| q156
    q156 --> q158
    q157 --> q145
    q157 --> q147
    q157 --> q149
    q157 --> q151
    q157 --> q153
    q158 --> q165
    q159 -->|"tok(&quot;*&quot;)"| q160
    q160 --> q166
    q161 -->|"tok(&quot;+&quot;)"| q162
    q162 --> q166
    q163 -->|"tok(&quot;?&quot;)"| q164
    q164 --> q166
    q165 --> q159
    q165 --> q161
    q165 --> q163
    q165 --> q166
    q166 --> q25
```

## Keyword

```mermaid
flowchart TD
    q26(["Keyword__Start (26)<br/>RuleStart"])
    q27(["Keyword__Stop (27)<br/>RuleStop"])
    q167["Keyword_Value_StringLiteral (167)<br/>Basic<br/>"]
    q168["Keyword__Basic (168)<br/>Basic<br/>"]

    q26 --> q167
    q167 -->|"tok(StringLiteral)"| q168
    q168 --> q27
```

## Assignment

```mermaid
flowchart TD
    q28(["Assignment__Start (28)<br/>RuleStart"])
    q29(["Assignment__Stop (29)<br/>RuleStop"])
    q169["Assignment_Property_ID (169)<br/>Basic<br/>"]
    q170["Assignment_Operator_PlusEquals (170)<br/>Basic<br/>"]
    q171["Assignment__Basic_0 (171)<br/>Basic<br/>"]
    q172["Assignment_Operator_Equals (172)<br/>Basic<br/>"]
    q173["Assignment__Basic_1 (173)<br/>Basic<br/>"]
    q174["Assignment_Operator_QuestionEquals (174)<br/>Basic<br/>"]
    q175["Assignment__Basic_2 (175)<br/>Basic<br/>"]
    q176["Assignment__Basic_3 (176)<br/>Basic<br/>"]
    q177["Assignment__BlockEnd (177)<br/>BlockEnd<br/>"]
    q178["Assignment__Basic_4 (178)<br/>Basic<br/>"]
    q179["Assignment__Basic_5 (179)<br/>Basic<br/>"]

    q28 --> q169
    q169 -->|"tok(ID)"| q176
    q170 -->|"tok(&quot;+=&quot;)"| q171
    q171 --> q177
    q172 -->|"tok(&quot;=&quot;)"| q173
    q173 --> q177
    q174 -->|"tok(&quot;?=&quot;)"| q175
    q175 --> q177
    q176 --> q170
    q176 --> q172
    q176 --> q174
    q177 --> q178
    q178 -.->|"[Assignable]"| q179
    q179 --> q29
```

## Assignable

```mermaid
flowchart TD
    q30(["Assignable__Start (30)<br/>RuleStart"])
    q31(["Assignable__Stop (31)<br/>RuleStop"])
    q180["Assignable__Basic_0 (180)<br/>Basic<br/>"]
    q181["Assignable__Basic_1 (181)<br/>Basic<br/>"]
    q182["Assignable__Basic_2 (182)<br/>Basic<br/>"]
    q183["Assignable__Basic_3 (183)<br/>Basic<br/>"]
    q184["Assignable__Basic_4 (184)<br/>Basic<br/>"]
    q185["Assignable__Basic_5 (185)<br/>Basic<br/>"]
    q186["Assignable_LeftParen (186)<br/>Basic<br/>"]
    q187["Assignable__Basic_6 (187)<br/>Basic<br/>"]
    q188["Assignable_RightParen (188)<br/>Basic<br/>"]
    q189["Assignable__Basic_7 (189)<br/>Basic<br/>"]
    q190["Assignable__Basic_8 (190)<br/>Basic<br/>"]
    q191["Assignable__BlockEnd (191)<br/>BlockEnd<br/>"]

    q30 --> q190
    q180 -.->|"[Keyword]"| q181
    q181 --> q191
    q182 -.->|"[RuleCall]"| q183
    q183 --> q191
    q184 -.->|"[CrossRef]"| q185
    q185 --> q191
    q186 -->|"tok(&quot;(&quot;)"| q187
    q187 -.->|"[AssignableAlternatives]"| q188
    q188 -->|"tok(&quot;)&quot;)"| q189
    q189 --> q191
    q190 --> q180
    q190 --> q182
    q190 --> q184
    q190 --> q186
    q191 --> q31
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q32(["AssignableWithoutAlts__Start (32)<br/>RuleStart"])
    q33(["AssignableWithoutAlts__Stop (33)<br/>RuleStop"])
    q192["AssignableWithoutAlts__Basic_0 (192)<br/>Basic<br/>"]
    q193["AssignableWithoutAlts__Basic_1 (193)<br/>Basic<br/>"]
    q194["AssignableWithoutAlts__Basic_2 (194)<br/>Basic<br/>"]
    q195["AssignableWithoutAlts__Basic_3 (195)<br/>Basic<br/>"]
    q196["AssignableWithoutAlts__Basic_4 (196)<br/>Basic<br/>"]
    q197["AssignableWithoutAlts__Basic_5 (197)<br/>Basic<br/>"]
    q198["AssignableWithoutAlts__Basic_6 (198)<br/>Basic<br/>"]
    q199["AssignableWithoutAlts__BlockEnd (199)<br/>BlockEnd<br/>"]

    q32 --> q198
    q192 -.->|"[Keyword]"| q193
    q193 --> q199
    q194 -.->|"[RuleCall]"| q195
    q195 --> q199
    q196 -.->|"[CrossRef]"| q197
    q197 --> q199
    q198 --> q192
    q198 --> q194
    q198 --> q196
    q199 --> q33
```

## AssignableAlternatives

```mermaid
flowchart TD
    q34(["AssignableAlternatives__Start (34)<br/>RuleStart"])
    q35(["AssignableAlternatives__Stop (35)<br/>RuleStop"])
    q200["AssignableAlternatives__Basic_0 (200)<br/>Basic<br/>"]
    q201["AssignableAlternatives_Pipe (201)<br/>Basic<br/>"]
    q202["AssignableAlternatives__Basic_1 (202)<br/>Basic<br/>"]
    q203["AssignableAlternatives__Basic_2 (203)<br/>Basic<br/>"]
    q204{"AssignableAlternatives__LoopBack (204)<br/>LoopBack<br/><br/>dec=5"}
    q205["AssignableAlternatives__LoopEnd (205)<br/>LoopEnd<br/>"]

    q34 --> q200
    q200 -.->|"[AssignableWithoutAlts]"| q201
    q201 -->|"tok(&quot;|&quot;)"| q202
    q201 --> q205
    q202 -.->|"[AssignableWithoutAlts]"| q203
    q203 --> q204
    q204 --> q201
    q204 --> q205
    q205 --> q35
```

## CrossRef

```mermaid
flowchart TD
    q36(["CrossRef__Start (36)<br/>RuleStart"])
    q37(["CrossRef__Stop (37)<br/>RuleStop"])
    q206["CrossRef_LeftBracket (206)<br/>Basic<br/>"]
    q207["CrossRef_Type_ID (207)<br/>Basic<br/>"]
    q208["CrossRef_Colon (208)<br/>Basic<br/>"]
    q209["CrossRef__Basic_0 (209)<br/>Basic<br/>"]
    q210["CrossRef__Basic_1 (210)<br/>Basic<br/>"]
    q211["CrossRef_RightBracket (211)<br/>Basic<br/>"]
    q212["CrossRef__Basic_2 (212)<br/>Basic<br/>"]

    q36 --> q206
    q206 -->|"tok(&quot;[&quot;)"| q207
    q207 -->|"tok(ID)"| q208
    q208 -->|"tok(&quot;:&quot;)"| q209
    q208 --> q210
    q209 -.->|"[RuleCall]"| q210
    q210 --> q211
    q211 -->|"tok(&quot;]&quot;)"| q212
    q212 --> q37
```

## RuleCall

```mermaid
flowchart TD
    q38(["RuleCall__Start (38)<br/>RuleStart"])
    q39(["RuleCall__Stop (39)<br/>RuleStop"])
    q213["RuleCall_Rule_ID (213)<br/>Basic<br/>"]
    q214["RuleCall__Basic (214)<br/>Basic<br/>"]

    q38 --> q213
    q213 -->|"tok(ID)"| q214
    q214 --> q39
```

## Action

```mermaid
flowchart TD
    q40(["Action__Start (40)<br/>RuleStart"])
    q41(["Action__Stop (41)<br/>RuleStop"])
    q215["Action_LeftBrace (215)<br/>Basic<br/>"]
    q216["Action_Type_ID (216)<br/>Basic<br/>"]
    q217["Action_Dot (217)<br/>Basic<br/>"]
    q218["Action_Property_ID (218)<br/>Basic<br/>"]
    q219["Action_Operator_PlusEquals (219)<br/>Basic<br/>"]
    q220["Action__Basic_0 (220)<br/>Basic<br/>"]
    q221["Action_Operator_Equals (221)<br/>Basic<br/>"]
    q222["Action__Basic_1 (222)<br/>Basic<br/>"]
    q223["Action__Basic_2 (223)<br/>Basic<br/>"]
    q224["Action__BlockEnd (224)<br/>BlockEnd<br/>"]
    q225["Action_current (225)<br/>Basic<br/>"]
    q226["Action__Basic_3 (226)<br/>Basic<br/>"]
    q227["Action_RightBrace (227)<br/>Basic<br/>"]
    q228["Action__Basic_4 (228)<br/>Basic<br/>"]

    q40 --> q215
    q215 -->|"tok(&quot;{&quot;)"| q216
    q216 -->|"tok(ID)"| q217
    q217 -->|"tok(&quot;.&quot;)"| q218
    q217 --> q226
    q218 -->|"tok(ID)"| q223
    q219 -->|"tok(&quot;+=&quot;)"| q220
    q220 --> q224
    q221 -->|"tok(&quot;=&quot;)"| q222
    q222 --> q224
    q223 --> q219
    q223 --> q221
    q224 --> q225
    q225 -->|"tok(&quot;current&quot;)"| q226
    q226 --> q227
    q227 -->|"tok(&quot;}&quot;)"| q228
    q228 --> q41
```

## CompositeRule

```mermaid
flowchart TD
    q42(["CompositeRule__Start (42)<br/>RuleStart"])
    q43(["CompositeRule__Stop (43)<br/>RuleStop"])
    q229["CompositeRule_composite (229)<br/>Basic<br/>"]
    q230["CompositeRule_Name_ID (230)<br/>Basic<br/>"]
    q231["CompositeRule_Colon (231)<br/>Basic<br/>"]
    q232["CompositeRule__Basic_0 (232)<br/>Basic<br/>"]
    q233["CompositeRule_Semicolon (233)<br/>Basic<br/>"]
    q234["CompositeRule__Basic_1 (234)<br/>Basic<br/>"]

    q42 --> q229
    q229 -->|"tok(&quot;composite&quot;)"| q230
    q230 -->|"tok(ID)"| q231
    q231 -->|"tok(&quot;:&quot;)"| q232
    q232 -.->|"[CompositeAlternatives]"| q233
    q233 -->|"tok(&quot;;&quot;)"| q234
    q234 --> q43
```

## CompositeAlternatives

```mermaid
flowchart TD
    q44(["CompositeAlternatives__Start (44)<br/>RuleStart"])
    q45(["CompositeAlternatives__Stop (45)<br/>RuleStop"])
    q235["CompositeAlternatives__Basic_0 (235)<br/>Basic<br/>"]
    q236["CompositeAlternatives_Pipe (236)<br/>Basic<br/>"]
    q237["CompositeAlternatives__Basic_1 (237)<br/>Basic<br/>"]
    q238["CompositeAlternatives__Basic_2 (238)<br/>Basic<br/>"]
    q239{"CompositeAlternatives__LoopBack (239)<br/>LoopBack<br/><br/>dec=6"}
    q240["CompositeAlternatives__LoopEnd (240)<br/>LoopEnd<br/>"]

    q44 --> q235
    q235 -.->|"[CompositeGroup]"| q236
    q236 -->|"tok(&quot;|&quot;)"| q237
    q236 --> q240
    q237 -.->|"[CompositeGroup]"| q238
    q238 --> q239
    q239 --> q236
    q239 --> q240
    q240 --> q45
```

## CompositeGroup

```mermaid
flowchart TD
    q46(["CompositeGroup__Start (46)<br/>RuleStart"])
    q47(["CompositeGroup__Stop (47)<br/>RuleStop"])
    q241["CompositeGroup__Basic_0 (241)<br/>Basic<br/>"]
    q242["CompositeGroup__Basic_1 (242)<br/>Basic<br/>"]
    q243["CompositeGroup__Basic_2 (243)<br/>Basic<br/>"]
    q244{"CompositeGroup__LoopBack (244)<br/>LoopBack<br/><br/>dec=7"}
    q245["CompositeGroup__LoopEnd (245)<br/>LoopEnd<br/>"]

    q46 --> q241
    q241 -.->|"[CompositeElement]"| q242
    q242 -.->|"[CompositeElement]"| q243
    q242 --> q245
    q243 --> q244
    q244 --> q242
    q244 --> q245
    q245 --> q47
```

## CompositeElement

```mermaid
flowchart TD
    q48(["CompositeElement__Start (48)<br/>RuleStart"])
    q49(["CompositeElement__Stop (49)<br/>RuleStop"])
    q246["CompositeElement__Basic_0 (246)<br/>Basic<br/>"]
    q247["CompositeElement__Basic_1 (247)<br/>Basic<br/>"]
    q248["CompositeElement__Basic_2 (248)<br/>Basic<br/>"]
    q249["CompositeElement__Basic_3 (249)<br/>Basic<br/>"]
    q250["CompositeElement_LeftParen (250)<br/>Basic<br/>"]
    q251["CompositeElement__Basic_4 (251)<br/>Basic<br/>"]
    q252["CompositeElement_RightParen (252)<br/>Basic<br/>"]
    q253["CompositeElement__Basic_5 (253)<br/>Basic<br/>"]
    q254["CompositeElement__Basic_6 (254)<br/>Basic<br/>"]
    q255["CompositeElement__BlockEnd_0 (255)<br/>BlockEnd<br/>"]
    q256["CompositeElement_Cardinality_Asterisk (256)<br/>Basic<br/>"]
    q257["CompositeElement__Basic_7 (257)<br/>Basic<br/>"]
    q258["CompositeElement_Cardinality_Plus (258)<br/>Basic<br/>"]
    q259["CompositeElement__Basic_8 (259)<br/>Basic<br/>"]
    q260["CompositeElement_Cardinality_Question (260)<br/>Basic<br/>"]
    q261["CompositeElement__Basic_9 (261)<br/>Basic<br/>"]
    q262["CompositeElement__Basic_10 (262)<br/>Basic<br/>"]
    q263["CompositeElement__BlockEnd_1 (263)<br/>BlockEnd<br/>"]

    q48 --> q254
    q246 -.->|"[Keyword]"| q247
    q247 --> q255
    q248 -.->|"[RuleCall]"| q249
    q249 --> q255
    q250 -->|"tok(&quot;(&quot;)"| q251
    q251 -.->|"[CompositeAlternatives]"| q252
    q252 -->|"tok(&quot;)&quot;)"| q253
    q253 --> q255
    q254 --> q246
    q254 --> q248
    q254 --> q250
    q255 --> q262
    q256 -->|"tok(&quot;*&quot;)"| q257
    q257 --> q263
    q258 -->|"tok(&quot;+&quot;)"| q259
    q259 --> q263
    q260 -->|"tok(&quot;?&quot;)"| q261
    q261 --> q263
    q262 --> q256
    q262 --> q258
    q262 --> q260
    q262 --> q263
    q263 --> q49
```

