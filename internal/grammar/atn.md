# Runtime ATN for grammar

## Grammar

```mermaid
flowchart TD
    q0(["Grammar__Start (0)<br/>RuleStart"])
    q1(["Grammar__Stop (1)<br/>RuleStop"])
    q52["Grammar_grammar (52)<br/>Basic<br/>"]
    q53["Grammar_Name_ID (53)<br/>Basic<br/>"]
    q54["Grammar_Semicolon (54)<br/>Basic<br/>"]
    q55["Grammar__Basic_0 (55)<br/>Basic<br/>"]
    q56["Grammar__Basic_1 (56)<br/>Basic<br/>"]
    q57["Grammar__Basic_2 (57)<br/>Basic<br/>"]
    q58["Grammar__Basic_3 (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_4 (59)<br/>Basic<br/>"]
    q60["Grammar__Basic_5 (60)<br/>Basic<br/>"]
    q61["Grammar__Basic_6 (61)<br/>Basic<br/>"]
    q62["Grammar__Basic_7 (62)<br/>Basic<br/>"]
    q63["Grammar__Basic_8 (63)<br/>Basic<br/>"]
    q64["Grammar__Basic_9 (64)<br/>Basic<br/>"]
    q65["Grammar__Basic_10 (65)<br/>Basic<br/>"]
    q66["Grammar__BlockEnd (66)<br/>BlockEnd<br/>"]
    q67{"Grammar__LoopEntry (67)<br/>LoopEntry<br/><br/>dec=0"}
    q68["Grammar__LoopEnd (68)<br/>LoopEnd<br/>"]
    q69["Grammar__LoopBack (69)<br/>LoopBack<br/>"]

    q0 --> q52
    q52 -->|"tok(&quot;grammar&quot;)"| q53
    q53 -->|"tok(ID)"| q54
    q54 -->|"tok(&quot;;&quot;)"| q67
    q55 -.->|"[ParserRule]"| q56
    q56 --> q66
    q57 -.->|"[Token]"| q58
    q58 --> q66
    q59 -.->|"[TokenGroup]"| q60
    q60 --> q66
    q61 -.->|"[Interface]"| q62
    q62 --> q66
    q63 -.->|"[CompositeRule]"| q64
    q64 --> q66
    q65 --> q55
    q65 --> q57
    q65 --> q59
    q65 --> q61
    q65 --> q63
    q66 --> q69
    q67 --> q65
    q67 --> q68
    q68 --> q1
    q69 --> q67
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q70["Interface_interface (70)<br/>Basic<br/>"]
    q71["Interface_Name_ID (71)<br/>Basic<br/>"]
    q72["Interface_extends (72)<br/>Basic<br/>"]
    q73["Interface_Extends_ID_0 (73)<br/>Basic<br/>"]
    q74["Interface_Comma (74)<br/>Basic<br/>"]
    q75["Interface_Extends_ID_1 (75)<br/>Basic<br/>"]
    q76["Interface__Basic_0 (76)<br/>Basic<br/>"]
    q77{"Interface__LoopEntry_0 (77)<br/>LoopEntry<br/><br/>dec=1"}
    q78["Interface__LoopEnd_0 (78)<br/>LoopEnd<br/>"]
    q79["Interface__LoopBack_0 (79)<br/>LoopBack<br/>"]
    q80["Interface_LeftBrace (80)<br/>Basic<br/>"]
    q81["Interface__Basic_1 (81)<br/>Basic<br/>"]
    q82["Interface__Basic_2 (82)<br/>Basic<br/>"]
    q83{"Interface__LoopEntry_1 (83)<br/>LoopEntry<br/><br/>dec=2"}
    q84["Interface__LoopEnd_1 (84)<br/>LoopEnd<br/>"]
    q85["Interface__LoopBack_1 (85)<br/>LoopBack<br/>"]
    q86["Interface_RightBrace (86)<br/>Basic<br/>"]
    q87["Interface__Basic_3 (87)<br/>Basic<br/>"]

    q2 --> q70
    q70 -->|"tok(&quot;interface&quot;)"| q71
    q71 -->|"tok(ID)"| q72
    q72 -->|"tok(&quot;extends&quot;)"| q73
    q72 --> q78
    q73 -->|"tok(ID)"| q77
    q74 -->|"tok(&quot;,&quot;)"| q75
    q75 -->|"tok(ID)"| q76
    q76 --> q79
    q77 --> q74
    q77 --> q78
    q78 --> q80
    q79 --> q77
    q80 -->|"tok(&quot;{&quot;)"| q83
    q81 -.->|"[Field]"| q82
    q82 --> q85
    q83 --> q81
    q83 --> q84
    q84 --> q86
    q85 --> q83
    q86 -->|"tok(&quot;}&quot;)"| q87
    q87 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q88["Field_Name_ID (88)<br/>Basic<br/>"]
    q89["Field__Basic_0 (89)<br/>Basic<br/>"]
    q90["Field__Basic_1 (90)<br/>Basic<br/>"]

    q4 --> q88
    q88 -->|"tok(ID)"| q89
    q89 -.->|"[FieldType]"| q90
    q90 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q91["FieldType__Basic_0 (91)<br/>Basic<br/>"]
    q92["FieldType__Basic_1 (92)<br/>Basic<br/>"]
    q93["FieldType__Basic_2 (93)<br/>Basic<br/>"]
    q94["FieldType__Basic_3 (94)<br/>Basic<br/>"]
    q95["FieldType__Basic_4 (95)<br/>Basic<br/>"]
    q96["FieldType__Basic_5 (96)<br/>Basic<br/>"]
    q97["FieldType__Basic_6 (97)<br/>Basic<br/>"]
    q98["FieldType__Basic_7 (98)<br/>Basic<br/>"]
    q99["FieldType__Basic_8 (99)<br/>Basic<br/>"]
    q100["FieldType__BlockEnd (100)<br/>BlockEnd<br/>"]

    q6 --> q99
    q91 -.->|"[SimpleType]"| q92
    q92 --> q100
    q93 -.->|"[ReferenceType]"| q94
    q94 --> q100
    q95 -.->|"[ArrayType]"| q96
    q96 --> q100
    q97 -.->|"[PrimitiveType]"| q98
    q98 --> q100
    q99 --> q91
    q99 --> q93
    q99 --> q95
    q99 --> q97
    q100 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q101["ArrayType_LeftBracket (101)<br/>Basic<br/>"]
    q102["ArrayType_RightBracket (102)<br/>Basic<br/>"]
    q103["ArrayType__Basic_0 (103)<br/>Basic<br/>"]
    q104["ArrayType__Basic_1 (104)<br/>Basic<br/>"]

    q8 --> q101
    q101 -->|"tok(&quot;[&quot;)"| q102
    q102 -->|"tok(&quot;]&quot;)"| q103
    q103 -.->|"[FieldType]"| q104
    q104 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q105["ReferenceType_Asterisk (105)<br/>Basic<br/>"]
    q106["ReferenceType_Type_ID (106)<br/>Basic<br/>"]
    q107["ReferenceType__Basic (107)<br/>Basic<br/>"]

    q10 --> q105
    q105 -->|"tok(&quot;*&quot;)"| q106
    q106 -->|"tok(ID)"| q107
    q107 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q108["SimpleType_Type_ID (108)<br/>Basic<br/>"]
    q109["SimpleType__Basic (109)<br/>Basic<br/>"]

    q12 --> q108
    q108 -->|"tok(ID)"| q109
    q109 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q110["PrimitiveType_Type_string (110)<br/>Basic<br/>"]
    q111["PrimitiveType__Basic_0 (111)<br/>Basic<br/>"]
    q112["PrimitiveType_Type_bool (112)<br/>Basic<br/>"]
    q113["PrimitiveType__Basic_1 (113)<br/>Basic<br/>"]
    q114["PrimitiveType_Type_composite (114)<br/>Basic<br/>"]
    q115["PrimitiveType__Basic_2 (115)<br/>Basic<br/>"]
    q116["PrimitiveType__Basic_3 (116)<br/>Basic<br/>"]
    q117["PrimitiveType__BlockEnd (117)<br/>BlockEnd<br/>"]

    q14 --> q116
    q110 -->|"tok(&quot;string&quot;)"| q111
    q111 --> q117
    q112 -->|"tok(&quot;bool&quot;)"| q113
    q113 --> q117
    q114 -->|"tok(&quot;composite&quot;)"| q115
    q115 --> q117
    q116 --> q110
    q116 --> q112
    q116 --> q114
    q117 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q118["ParserRule_Name_ID (118)<br/>Basic<br/>"]
    q119["ParserRule_returns (119)<br/>Basic<br/>"]
    q120["ParserRule_ReturnType_ID (120)<br/>Basic<br/>"]
    q121["ParserRule__Basic_0 (121)<br/>Basic<br/>"]
    q122["ParserRule_Colon (122)<br/>Basic<br/>"]
    q123["ParserRule__Basic_1 (123)<br/>Basic<br/>"]
    q124["ParserRule_Semicolon (124)<br/>Basic<br/>"]
    q125["ParserRule__Basic_2 (125)<br/>Basic<br/>"]

    q16 --> q118
    q118 -->|"tok(ID)"| q119
    q119 -->|"tok(&quot;returns&quot;)"| q120
    q119 --> q121
    q120 -->|"tok(ID)"| q121
    q121 --> q122
    q122 -->|"tok(&quot;:&quot;)"| q123
    q123 -.->|"[Alternatives]"| q124
    q124 -->|"tok(&quot;;&quot;)"| q125
    q125 --> q17
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q126["Token_Type_hidden (126)<br/>Basic<br/>"]
    q127["Token__Basic_0 (127)<br/>Basic<br/>"]
    q128["Token_Type_comment (128)<br/>Basic<br/>"]
    q129["Token__Basic_1 (129)<br/>Basic<br/>"]
    q130["Token__Basic_2 (130)<br/>Basic<br/>"]
    q131["Token__BlockEnd (131)<br/>BlockEnd<br/>"]
    q132["Token_token (132)<br/>Basic<br/>"]
    q133["Token_Name_ID (133)<br/>Basic<br/>"]
    q134["Token_Colon (134)<br/>Basic<br/>"]
    q135["Token_Regexp_RegexLiteral (135)<br/>Basic<br/>"]
    q136["Token_Semicolon (136)<br/>Basic<br/>"]
    q137["Token__Basic_3 (137)<br/>Basic<br/>"]

    q18 --> q130
    q126 -->|"tok(&quot;hidden&quot;)"| q127
    q127 --> q131
    q128 -->|"tok(&quot;comment&quot;)"| q129
    q129 --> q131
    q130 --> q126
    q130 --> q128
    q130 --> q131
    q131 --> q132
    q132 -->|"tok(&quot;token&quot;)"| q133
    q133 -->|"tok(ID)"| q134
    q134 -->|"tok(&quot;:&quot;)"| q135
    q135 -->|"tok(RegexLiteral)"| q136
    q136 -->|"tok(&quot;;&quot;)"| q137
    q137 --> q19
```

## TokenGroup

```mermaid
flowchart TD
    q20(["TokenGroup__Start (20)<br/>RuleStart"])
    q21(["TokenGroup__Stop (21)<br/>RuleStop"])
    q138["TokenGroup_token (138)<br/>Basic<br/>"]
    q139["TokenGroup_group (139)<br/>Basic<br/>"]
    q140["TokenGroup_Name_ID (140)<br/>Basic<br/>"]
    q141["TokenGroup_LeftBrace (141)<br/>Basic<br/>"]
    q142["TokenGroup_TokenRefs_ID (142)<br/>Basic<br/>"]
    q143["TokenGroup__Basic_0 (143)<br/>Basic<br/>"]
    q144["TokenGroup_keywords (144)<br/>Basic<br/>"]
    q145["TokenGroup_Regexps_RegexLiteral (145)<br/>Basic<br/>"]
    q146["TokenGroup__Basic_1 (146)<br/>Basic<br/>"]
    q147["TokenGroup__Basic_2 (147)<br/>Basic<br/>"]
    q148["TokenGroup__Basic_3 (148)<br/>Basic<br/>"]
    q149["TokenGroup__Basic_4 (149)<br/>Basic<br/>"]
    q150["TokenGroup__BlockEnd (150)<br/>BlockEnd<br/>"]
    q151{"TokenGroup__LoopEntry (151)<br/>LoopEntry<br/><br/>dec=3"}
    q152["TokenGroup__LoopEnd (152)<br/>LoopEnd<br/>"]
    q153["TokenGroup__LoopBack (153)<br/>LoopBack<br/>"]
    q154["TokenGroup_RightBrace (154)<br/>Basic<br/>"]
    q155["TokenGroup__Basic_5 (155)<br/>Basic<br/>"]

    q20 --> q138
    q138 -->|"tok(&quot;token&quot;)"| q139
    q139 -->|"tok(&quot;group&quot;)"| q140
    q140 -->|"tok(ID)"| q141
    q141 -->|"tok(&quot;{&quot;)"| q151
    q142 -->|"tok(ID)"| q143
    q143 --> q150
    q144 -->|"tok(&quot;keywords&quot;)"| q145
    q145 -->|"tok(RegexLiteral)"| q146
    q146 --> q150
    q147 -.->|"[Keyword]"| q148
    q148 --> q150
    q149 --> q142
    q149 --> q144
    q149 --> q147
    q150 --> q153
    q151 --> q149
    q151 --> q152
    q152 --> q154
    q153 --> q151
    q154 -->|"tok(&quot;}&quot;)"| q155
    q155 --> q21
```

## Alternatives

```mermaid
flowchart TD
    q22(["Alternatives__Start (22)<br/>RuleStart"])
    q23(["Alternatives__Stop (23)<br/>RuleStop"])
    q156["Alternatives__Basic_0 (156)<br/>Basic<br/>"]
    q157["Alternatives_Pipe (157)<br/>Basic<br/>"]
    q158["Alternatives__Basic_1 (158)<br/>Basic<br/>"]
    q159["Alternatives__Basic_2 (159)<br/>Basic<br/>"]
    q160{"Alternatives__LoopBack (160)<br/>LoopBack<br/><br/>dec=4"}
    q161["Alternatives__LoopEnd (161)<br/>LoopEnd<br/>"]

    q22 --> q156
    q156 -.->|"[Group]"| q157
    q157 -->|"tok(&quot;|&quot;)"| q158
    q157 --> q161
    q158 -.->|"[Group]"| q159
    q159 --> q160
    q160 --> q157
    q160 --> q161
    q161 --> q23
```

## Group

```mermaid
flowchart TD
    q24(["Group__Start (24)<br/>RuleStart"])
    q25(["Group__Stop (25)<br/>RuleStop"])
    q162["Group__Basic_0 (162)<br/>Basic<br/>"]
    q163["Group__Basic_1 (163)<br/>Basic<br/>"]
    q164["Group__Basic_2 (164)<br/>Basic<br/>"]
    q165{"Group__LoopBack (165)<br/>LoopBack<br/><br/>dec=5"}
    q166["Group__LoopEnd (166)<br/>LoopEnd<br/>"]

    q24 --> q162
    q162 -.->|"[Element]"| q163
    q163 -.->|"[Element]"| q164
    q163 --> q166
    q164 --> q165
    q165 --> q163
    q165 --> q166
    q166 --> q25
```

## Element

```mermaid
flowchart TD
    q26(["Element__Start (26)<br/>RuleStart"])
    q27(["Element__Stop (27)<br/>RuleStop"])
    q167["Element__Basic_0 (167)<br/>Basic<br/>"]
    q168["Element__Basic_1 (168)<br/>Basic<br/>"]
    q169["Element__Basic_2 (169)<br/>Basic<br/>"]
    q170["Element__Basic_3 (170)<br/>Basic<br/>"]
    q171["Element__Basic_4 (171)<br/>Basic<br/>"]
    q172["Element__Basic_5 (172)<br/>Basic<br/>"]
    q173["Element__Basic_6 (173)<br/>Basic<br/>"]
    q174["Element__Basic_7 (174)<br/>Basic<br/>"]
    q175["Element_LeftParen (175)<br/>Basic<br/>"]
    q176["Element__Basic_8 (176)<br/>Basic<br/>"]
    q177["Element_RightParen (177)<br/>Basic<br/>"]
    q178["Element__Basic_9 (178)<br/>Basic<br/>"]
    q179["Element__Basic_10 (179)<br/>Basic<br/>"]
    q180["Element__BlockEnd_0 (180)<br/>BlockEnd<br/>"]
    q181["Element_Cardinality_Asterisk (181)<br/>Basic<br/>"]
    q182["Element__Basic_11 (182)<br/>Basic<br/>"]
    q183["Element_Cardinality_Plus (183)<br/>Basic<br/>"]
    q184["Element__Basic_12 (184)<br/>Basic<br/>"]
    q185["Element_Cardinality_Question (185)<br/>Basic<br/>"]
    q186["Element__Basic_13 (186)<br/>Basic<br/>"]
    q187["Element__Basic_14 (187)<br/>Basic<br/>"]
    q188["Element__BlockEnd_1 (188)<br/>BlockEnd<br/>"]

    q26 --> q179
    q167 -.->|"[Keyword]"| q168
    q168 --> q180
    q169 -.->|"[Assignment]"| q170
    q170 --> q180
    q171 -.->|"[RuleCall]"| q172
    q172 --> q180
    q173 -.->|"[Action]"| q174
    q174 --> q180
    q175 -->|"tok(&quot;(&quot;)"| q176
    q176 -.->|"[Alternatives]"| q177
    q177 -->|"tok(&quot;)&quot;)"| q178
    q178 --> q180
    q179 --> q167
    q179 --> q169
    q179 --> q171
    q179 --> q173
    q179 --> q175
    q180 --> q187
    q181 -->|"tok(&quot;*&quot;)"| q182
    q182 --> q188
    q183 -->|"tok(&quot;+&quot;)"| q184
    q184 --> q188
    q185 -->|"tok(&quot;?&quot;)"| q186
    q186 --> q188
    q187 --> q181
    q187 --> q183
    q187 --> q185
    q187 --> q188
    q188 --> q27
```

## Keyword

```mermaid
flowchart TD
    q28(["Keyword__Start (28)<br/>RuleStart"])
    q29(["Keyword__Stop (29)<br/>RuleStop"])
    q189["Keyword_Value_StringLiteral (189)<br/>Basic<br/>"]
    q190["Keyword__Basic (190)<br/>Basic<br/>"]

    q28 --> q189
    q189 -->|"tok(StringLiteral)"| q190
    q190 --> q29
```

## Assignment

```mermaid
flowchart TD
    q30(["Assignment__Start (30)<br/>RuleStart"])
    q31(["Assignment__Stop (31)<br/>RuleStop"])
    q191["Assignment_Property_ID (191)<br/>Basic<br/>"]
    q192["Assignment_Operator_PlusEquals (192)<br/>Basic<br/>"]
    q193["Assignment__Basic_0 (193)<br/>Basic<br/>"]
    q194["Assignment_Operator_Equals (194)<br/>Basic<br/>"]
    q195["Assignment__Basic_1 (195)<br/>Basic<br/>"]
    q196["Assignment_Operator_QuestionEquals (196)<br/>Basic<br/>"]
    q197["Assignment__Basic_2 (197)<br/>Basic<br/>"]
    q198["Assignment__Basic_3 (198)<br/>Basic<br/>"]
    q199["Assignment__BlockEnd (199)<br/>BlockEnd<br/>"]
    q200["Assignment__Basic_4 (200)<br/>Basic<br/>"]
    q201["Assignment__Basic_5 (201)<br/>Basic<br/>"]

    q30 --> q191
    q191 -->|"tok(ID)"| q198
    q192 -->|"tok(&quot;+=&quot;)"| q193
    q193 --> q199
    q194 -->|"tok(&quot;=&quot;)"| q195
    q195 --> q199
    q196 -->|"tok(&quot;?=&quot;)"| q197
    q197 --> q199
    q198 --> q192
    q198 --> q194
    q198 --> q196
    q199 --> q200
    q200 -.->|"[Assignable]"| q201
    q201 --> q31
```

## Assignable

```mermaid
flowchart TD
    q32(["Assignable__Start (32)<br/>RuleStart"])
    q33(["Assignable__Stop (33)<br/>RuleStop"])
    q202["Assignable__Basic_0 (202)<br/>Basic<br/>"]
    q203["Assignable__Basic_1 (203)<br/>Basic<br/>"]
    q204["Assignable__Basic_2 (204)<br/>Basic<br/>"]
    q205["Assignable__Basic_3 (205)<br/>Basic<br/>"]
    q206["Assignable__Basic_4 (206)<br/>Basic<br/>"]
    q207["Assignable__Basic_5 (207)<br/>Basic<br/>"]
    q208["Assignable_LeftParen (208)<br/>Basic<br/>"]
    q209["Assignable__Basic_6 (209)<br/>Basic<br/>"]
    q210["Assignable_RightParen (210)<br/>Basic<br/>"]
    q211["Assignable__Basic_7 (211)<br/>Basic<br/>"]
    q212["Assignable__Basic_8 (212)<br/>Basic<br/>"]
    q213["Assignable__BlockEnd (213)<br/>BlockEnd<br/>"]

    q32 --> q212
    q202 -.->|"[Keyword]"| q203
    q203 --> q213
    q204 -.->|"[RuleCall]"| q205
    q205 --> q213
    q206 -.->|"[CrossRef]"| q207
    q207 --> q213
    q208 -->|"tok(&quot;(&quot;)"| q209
    q209 -.->|"[AssignableAlternatives]"| q210
    q210 -->|"tok(&quot;)&quot;)"| q211
    q211 --> q213
    q212 --> q202
    q212 --> q204
    q212 --> q206
    q212 --> q208
    q213 --> q33
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q34(["AssignableWithoutAlts__Start (34)<br/>RuleStart"])
    q35(["AssignableWithoutAlts__Stop (35)<br/>RuleStop"])
    q214["AssignableWithoutAlts__Basic_0 (214)<br/>Basic<br/>"]
    q215["AssignableWithoutAlts__Basic_1 (215)<br/>Basic<br/>"]
    q216["AssignableWithoutAlts__Basic_2 (216)<br/>Basic<br/>"]
    q217["AssignableWithoutAlts__Basic_3 (217)<br/>Basic<br/>"]
    q218["AssignableWithoutAlts__Basic_4 (218)<br/>Basic<br/>"]
    q219["AssignableWithoutAlts__Basic_5 (219)<br/>Basic<br/>"]
    q220["AssignableWithoutAlts__Basic_6 (220)<br/>Basic<br/>"]
    q221["AssignableWithoutAlts__BlockEnd (221)<br/>BlockEnd<br/>"]

    q34 --> q220
    q214 -.->|"[Keyword]"| q215
    q215 --> q221
    q216 -.->|"[RuleCall]"| q217
    q217 --> q221
    q218 -.->|"[CrossRef]"| q219
    q219 --> q221
    q220 --> q214
    q220 --> q216
    q220 --> q218
    q221 --> q35
```

## AssignableAlternatives

```mermaid
flowchart TD
    q36(["AssignableAlternatives__Start (36)<br/>RuleStart"])
    q37(["AssignableAlternatives__Stop (37)<br/>RuleStop"])
    q222["AssignableAlternatives__Basic_0 (222)<br/>Basic<br/>"]
    q223["AssignableAlternatives_Pipe (223)<br/>Basic<br/>"]
    q224["AssignableAlternatives__Basic_1 (224)<br/>Basic<br/>"]
    q225["AssignableAlternatives__Basic_2 (225)<br/>Basic<br/>"]
    q226{"AssignableAlternatives__LoopBack (226)<br/>LoopBack<br/><br/>dec=6"}
    q227["AssignableAlternatives__LoopEnd (227)<br/>LoopEnd<br/>"]

    q36 --> q222
    q222 -.->|"[AssignableWithoutAlts]"| q223
    q223 -->|"tok(&quot;|&quot;)"| q224
    q223 --> q227
    q224 -.->|"[AssignableWithoutAlts]"| q225
    q225 --> q226
    q226 --> q223
    q226 --> q227
    q227 --> q37
```

## CrossRef

```mermaid
flowchart TD
    q38(["CrossRef__Start (38)<br/>RuleStart"])
    q39(["CrossRef__Stop (39)<br/>RuleStop"])
    q228["CrossRef_LeftBracket (228)<br/>Basic<br/>"]
    q229["CrossRef_Type_ID (229)<br/>Basic<br/>"]
    q230["CrossRef_Colon (230)<br/>Basic<br/>"]
    q231["CrossRef__Basic_0 (231)<br/>Basic<br/>"]
    q232["CrossRef__Basic_1 (232)<br/>Basic<br/>"]
    q233["CrossRef_RightBracket (233)<br/>Basic<br/>"]
    q234["CrossRef__Basic_2 (234)<br/>Basic<br/>"]

    q38 --> q228
    q228 -->|"tok(&quot;[&quot;)"| q229
    q229 -->|"tok(ID)"| q230
    q230 -->|"tok(&quot;:&quot;)"| q231
    q230 --> q232
    q231 -.->|"[RuleCall]"| q232
    q232 --> q233
    q233 -->|"tok(&quot;]&quot;)"| q234
    q234 --> q39
```

## RuleCall

```mermaid
flowchart TD
    q40(["RuleCall__Start (40)<br/>RuleStart"])
    q41(["RuleCall__Stop (41)<br/>RuleStop"])
    q235["RuleCall_Rule_ID (235)<br/>Basic<br/>"]
    q236["RuleCall__Basic (236)<br/>Basic<br/>"]

    q40 --> q235
    q235 -->|"tok(ID)"| q236
    q236 --> q41
```

## Action

```mermaid
flowchart TD
    q42(["Action__Start (42)<br/>RuleStart"])
    q43(["Action__Stop (43)<br/>RuleStop"])
    q237["Action_LeftBrace (237)<br/>Basic<br/>"]
    q238["Action_Type_ID (238)<br/>Basic<br/>"]
    q239["Action_Dot (239)<br/>Basic<br/>"]
    q240["Action_Property_ID (240)<br/>Basic<br/>"]
    q241["Action_Operator_PlusEquals (241)<br/>Basic<br/>"]
    q242["Action__Basic_0 (242)<br/>Basic<br/>"]
    q243["Action_Operator_Equals (243)<br/>Basic<br/>"]
    q244["Action__Basic_1 (244)<br/>Basic<br/>"]
    q245["Action__Basic_2 (245)<br/>Basic<br/>"]
    q246["Action__BlockEnd (246)<br/>BlockEnd<br/>"]
    q247["Action_current (247)<br/>Basic<br/>"]
    q248["Action__Basic_3 (248)<br/>Basic<br/>"]
    q249["Action_RightBrace (249)<br/>Basic<br/>"]
    q250["Action__Basic_4 (250)<br/>Basic<br/>"]

    q42 --> q237
    q237 -->|"tok(&quot;{&quot;)"| q238
    q238 -->|"tok(ID)"| q239
    q239 -->|"tok(&quot;.&quot;)"| q240
    q239 --> q248
    q240 -->|"tok(ID)"| q245
    q241 -->|"tok(&quot;+=&quot;)"| q242
    q242 --> q246
    q243 -->|"tok(&quot;=&quot;)"| q244
    q244 --> q246
    q245 --> q241
    q245 --> q243
    q246 --> q247
    q247 -->|"tok(&quot;current&quot;)"| q248
    q248 --> q249
    q249 -->|"tok(&quot;}&quot;)"| q250
    q250 --> q43
```

## CompositeRule

```mermaid
flowchart TD
    q44(["CompositeRule__Start (44)<br/>RuleStart"])
    q45(["CompositeRule__Stop (45)<br/>RuleStop"])
    q251["CompositeRule_composite (251)<br/>Basic<br/>"]
    q252["CompositeRule_Name_ID (252)<br/>Basic<br/>"]
    q253["CompositeRule_Colon (253)<br/>Basic<br/>"]
    q254["CompositeRule__Basic_0 (254)<br/>Basic<br/>"]
    q255["CompositeRule_Semicolon (255)<br/>Basic<br/>"]
    q256["CompositeRule__Basic_1 (256)<br/>Basic<br/>"]

    q44 --> q251
    q251 -->|"tok(&quot;composite&quot;)"| q252
    q252 -->|"tok(ID)"| q253
    q253 -->|"tok(&quot;:&quot;)"| q254
    q254 -.->|"[CompositeAlternatives]"| q255
    q255 -->|"tok(&quot;;&quot;)"| q256
    q256 --> q45
```

## CompositeAlternatives

```mermaid
flowchart TD
    q46(["CompositeAlternatives__Start (46)<br/>RuleStart"])
    q47(["CompositeAlternatives__Stop (47)<br/>RuleStop"])
    q257["CompositeAlternatives__Basic_0 (257)<br/>Basic<br/>"]
    q258["CompositeAlternatives_Pipe (258)<br/>Basic<br/>"]
    q259["CompositeAlternatives__Basic_1 (259)<br/>Basic<br/>"]
    q260["CompositeAlternatives__Basic_2 (260)<br/>Basic<br/>"]
    q261{"CompositeAlternatives__LoopBack (261)<br/>LoopBack<br/><br/>dec=7"}
    q262["CompositeAlternatives__LoopEnd (262)<br/>LoopEnd<br/>"]

    q46 --> q257
    q257 -.->|"[CompositeGroup]"| q258
    q258 -->|"tok(&quot;|&quot;)"| q259
    q258 --> q262
    q259 -.->|"[CompositeGroup]"| q260
    q260 --> q261
    q261 --> q258
    q261 --> q262
    q262 --> q47
```

## CompositeGroup

```mermaid
flowchart TD
    q48(["CompositeGroup__Start (48)<br/>RuleStart"])
    q49(["CompositeGroup__Stop (49)<br/>RuleStop"])
    q263["CompositeGroup__Basic_0 (263)<br/>Basic<br/>"]
    q264["CompositeGroup__Basic_1 (264)<br/>Basic<br/>"]
    q265["CompositeGroup__Basic_2 (265)<br/>Basic<br/>"]
    q266{"CompositeGroup__LoopBack (266)<br/>LoopBack<br/><br/>dec=8"}
    q267["CompositeGroup__LoopEnd (267)<br/>LoopEnd<br/>"]

    q48 --> q263
    q263 -.->|"[CompositeElement]"| q264
    q264 -.->|"[CompositeElement]"| q265
    q264 --> q267
    q265 --> q266
    q266 --> q264
    q266 --> q267
    q267 --> q49
```

## CompositeElement

```mermaid
flowchart TD
    q50(["CompositeElement__Start (50)<br/>RuleStart"])
    q51(["CompositeElement__Stop (51)<br/>RuleStop"])
    q268["CompositeElement__Basic_0 (268)<br/>Basic<br/>"]
    q269["CompositeElement__Basic_1 (269)<br/>Basic<br/>"]
    q270["CompositeElement__Basic_2 (270)<br/>Basic<br/>"]
    q271["CompositeElement__Basic_3 (271)<br/>Basic<br/>"]
    q272["CompositeElement_LeftParen (272)<br/>Basic<br/>"]
    q273["CompositeElement__Basic_4 (273)<br/>Basic<br/>"]
    q274["CompositeElement_RightParen (274)<br/>Basic<br/>"]
    q275["CompositeElement__Basic_5 (275)<br/>Basic<br/>"]
    q276["CompositeElement__Basic_6 (276)<br/>Basic<br/>"]
    q277["CompositeElement__BlockEnd_0 (277)<br/>BlockEnd<br/>"]
    q278["CompositeElement_Cardinality_Asterisk (278)<br/>Basic<br/>"]
    q279["CompositeElement__Basic_7 (279)<br/>Basic<br/>"]
    q280["CompositeElement_Cardinality_Plus (280)<br/>Basic<br/>"]
    q281["CompositeElement__Basic_8 (281)<br/>Basic<br/>"]
    q282["CompositeElement_Cardinality_Question (282)<br/>Basic<br/>"]
    q283["CompositeElement__Basic_9 (283)<br/>Basic<br/>"]
    q284["CompositeElement__Basic_10 (284)<br/>Basic<br/>"]
    q285["CompositeElement__BlockEnd_1 (285)<br/>BlockEnd<br/>"]

    q50 --> q276
    q268 -.->|"[Keyword]"| q269
    q269 --> q277
    q270 -.->|"[RuleCall]"| q271
    q271 --> q277
    q272 -->|"tok(&quot;(&quot;)"| q273
    q273 -.->|"[CompositeAlternatives]"| q274
    q274 -->|"tok(&quot;)&quot;)"| q275
    q275 --> q277
    q276 --> q268
    q276 --> q270
    q276 --> q272
    q277 --> q284
    q278 -->|"tok(&quot;*&quot;)"| q279
    q279 --> q285
    q280 -->|"tok(&quot;+&quot;)"| q281
    q281 --> q285
    q282 -->|"tok(&quot;?&quot;)"| q283
    q283 --> q285
    q284 --> q278
    q284 --> q280
    q284 --> q282
    q284 --> q285
    q285 --> q51
```

