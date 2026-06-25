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
    q56{"Grammar__Basic_1 (56)<br/>Basic<br/><br/>dec=0"}
    q57["Grammar__Basic_2 (57)<br/>Basic<br/>"]
    q58["Grammar__Basic_3 (58)<br/>Basic<br/>"]
    q59["Grammar__Basic_4 (59)<br/>Basic<br/>"]
    q60["Grammar__Basic_5 (60)<br/>Basic<br/>"]
    q61["Grammar__Basic_6 (61)<br/>Basic<br/>"]
    q62["Grammar__Basic_7 (62)<br/>Basic<br/>"]
    q63["Grammar__Basic_8 (63)<br/>Basic<br/>"]
    q64["Grammar__Basic_9 (64)<br/>Basic<br/>"]
    q65["Grammar__Basic_10 (65)<br/>Basic<br/>"]
    q66["Grammar__Basic_11 (66)<br/>Basic<br/>"]
    q67{"Grammar__Basic_12 (67)<br/>Basic<br/><br/>dec=1"}
    q68["Grammar__BlockEnd (68)<br/>BlockEnd<br/>"]
    q69{"Grammar__LoopEntry (69)<br/>LoopEntry<br/><br/>dec=2"}
    q70["Grammar__LoopEnd (70)<br/>LoopEnd<br/>"]
    q71["Grammar__LoopBack (71)<br/>LoopBack<br/>"]

    q0 --> q52
    q52 -->|"tok(&quot;grammar&quot;)"| q53
    q53 -->|"tok(ID)"| q56
    q54 -->|"tok(&quot;;&quot;)"| q55
    q55 --> q69
    q56 --> q54
    q56 --> q55
    q57 -.->|"[ParserRule]"| q58
    q58 --> q68
    q59 -.->|"[Token]"| q60
    q60 --> q68
    q61 -.->|"[TokenGroup]"| q62
    q62 --> q68
    q63 -.->|"[Interface]"| q64
    q64 --> q68
    q65 -.->|"[CompositeRule]"| q66
    q66 --> q68
    q67 --> q57
    q67 --> q59
    q67 --> q61
    q67 --> q63
    q67 --> q65
    q68 --> q71
    q69 --> q67
    q69 --> q70
    q70 --> q1
    q71 --> q69
```

## Interface

```mermaid
flowchart TD
    q2(["Interface__Start (2)<br/>RuleStart"])
    q3(["Interface__Stop (3)<br/>RuleStop"])
    q72["Interface_interface (72)<br/>Basic<br/>"]
    q73["Interface_Name_ID (73)<br/>Basic<br/>"]
    q74["Interface_extends (74)<br/>Basic<br/>"]
    q75["Interface_Extends_ID_0 (75)<br/>Basic<br/>"]
    q76["Interface_Comma (76)<br/>Basic<br/>"]
    q77["Interface_Extends_ID_1 (77)<br/>Basic<br/>"]
    q78["Interface__Basic_0 (78)<br/>Basic<br/>"]
    q79{"Interface__LoopEntry_0 (79)<br/>LoopEntry<br/><br/>dec=3"}
    q80["Interface__LoopEnd_0 (80)<br/>LoopEnd<br/>"]
    q81["Interface__LoopBack_0 (81)<br/>LoopBack<br/>"]
    q82{"Interface__Basic_1 (82)<br/>Basic<br/><br/>dec=4"}
    q83["Interface_LeftBrace (83)<br/>Basic<br/>"]
    q84["Interface__Basic_2 (84)<br/>Basic<br/>"]
    q85["Interface__Basic_3 (85)<br/>Basic<br/>"]
    q86{"Interface__LoopEntry_1 (86)<br/>LoopEntry<br/><br/>dec=5"}
    q87["Interface__LoopEnd_1 (87)<br/>LoopEnd<br/>"]
    q88["Interface__LoopBack_1 (88)<br/>LoopBack<br/>"]
    q89["Interface_RightBrace (89)<br/>Basic<br/>"]
    q90["Interface__Basic_4 (90)<br/>Basic<br/>"]

    q2 --> q72
    q72 -->|"tok(&quot;interface&quot;)"| q73
    q73 -->|"tok(ID)"| q82
    q74 -->|"tok(&quot;extends&quot;)"| q75
    q75 -->|"tok(ID)"| q79
    q76 -->|"tok(&quot;,&quot;)"| q77
    q77 -->|"tok(ID)"| q78
    q78 --> q81
    q79 --> q76
    q79 --> q80
    q80 --> q83
    q81 --> q79
    q82 --> q74
    q82 --> q80
    q83 -->|"tok(&quot;{&quot;)"| q86
    q84 -.->|"[Field]"| q85
    q85 --> q88
    q86 --> q84
    q86 --> q87
    q87 --> q89
    q88 --> q86
    q89 -->|"tok(&quot;}&quot;)"| q90
    q90 --> q3
```

## Field

```mermaid
flowchart TD
    q4(["Field__Start (4)<br/>RuleStart"])
    q5(["Field__Stop (5)<br/>RuleStop"])
    q91["Field_Name_ID (91)<br/>Basic<br/>"]
    q92["Field__Basic_0 (92)<br/>Basic<br/>"]
    q93["Field__Basic_1 (93)<br/>Basic<br/>"]

    q4 --> q91
    q91 -->|"tok(ID)"| q92
    q92 -.->|"[FieldType]"| q93
    q93 --> q5
```

## FieldType

```mermaid
flowchart TD
    q6(["FieldType__Start (6)<br/>RuleStart"])
    q7(["FieldType__Stop (7)<br/>RuleStop"])
    q94["FieldType__Basic_0 (94)<br/>Basic<br/>"]
    q95["FieldType__Basic_1 (95)<br/>Basic<br/>"]
    q96["FieldType__Basic_2 (96)<br/>Basic<br/>"]
    q97["FieldType__Basic_3 (97)<br/>Basic<br/>"]
    q98["FieldType__Basic_4 (98)<br/>Basic<br/>"]
    q99["FieldType__Basic_5 (99)<br/>Basic<br/>"]
    q100["FieldType__Basic_6 (100)<br/>Basic<br/>"]
    q101["FieldType__Basic_7 (101)<br/>Basic<br/>"]
    q102{"FieldType__Basic_8 (102)<br/>Basic<br/><br/>dec=6"}
    q103["FieldType__BlockEnd (103)<br/>BlockEnd<br/>"]

    q6 --> q102
    q94 -.->|"[SimpleType]"| q95
    q95 --> q103
    q96 -.->|"[ReferenceType]"| q97
    q97 --> q103
    q98 -.->|"[ArrayType]"| q99
    q99 --> q103
    q100 -.->|"[PrimitiveType]"| q101
    q101 --> q103
    q102 --> q94
    q102 --> q96
    q102 --> q98
    q102 --> q100
    q103 --> q7
```

## ArrayType

```mermaid
flowchart TD
    q8(["ArrayType__Start (8)<br/>RuleStart"])
    q9(["ArrayType__Stop (9)<br/>RuleStop"])
    q104["ArrayType_LeftBracket (104)<br/>Basic<br/>"]
    q105["ArrayType_RightBracket (105)<br/>Basic<br/>"]
    q106["ArrayType__Basic_0 (106)<br/>Basic<br/>"]
    q107["ArrayType__Basic_1 (107)<br/>Basic<br/>"]

    q8 --> q104
    q104 -->|"tok(&quot;[&quot;)"| q105
    q105 -->|"tok(&quot;]&quot;)"| q106
    q106 -.->|"[FieldType]"| q107
    q107 --> q9
```

## ReferenceType

```mermaid
flowchart TD
    q10(["ReferenceType__Start (10)<br/>RuleStart"])
    q11(["ReferenceType__Stop (11)<br/>RuleStop"])
    q108["ReferenceType_Asterisk (108)<br/>Basic<br/>"]
    q109["ReferenceType_Type_ID (109)<br/>Basic<br/>"]
    q110["ReferenceType__Basic (110)<br/>Basic<br/>"]

    q10 --> q108
    q108 -->|"tok(&quot;*&quot;)"| q109
    q109 -->|"tok(ID)"| q110
    q110 --> q11
```

## SimpleType

```mermaid
flowchart TD
    q12(["SimpleType__Start (12)<br/>RuleStart"])
    q13(["SimpleType__Stop (13)<br/>RuleStop"])
    q111["SimpleType_Type_ID (111)<br/>Basic<br/>"]
    q112["SimpleType__Basic (112)<br/>Basic<br/>"]

    q12 --> q111
    q111 -->|"tok(ID)"| q112
    q112 --> q13
```

## PrimitiveType

```mermaid
flowchart TD
    q14(["PrimitiveType__Start (14)<br/>RuleStart"])
    q15(["PrimitiveType__Stop (15)<br/>RuleStop"])
    q113["PrimitiveType_Type_string (113)<br/>Basic<br/>"]
    q114["PrimitiveType__Basic_0 (114)<br/>Basic<br/>"]
    q115["PrimitiveType_Type_bool (115)<br/>Basic<br/>"]
    q116["PrimitiveType__Basic_1 (116)<br/>Basic<br/>"]
    q117["PrimitiveType_Type_composite (117)<br/>Basic<br/>"]
    q118["PrimitiveType__Basic_2 (118)<br/>Basic<br/>"]
    q119{"PrimitiveType__Basic_3 (119)<br/>Basic<br/><br/>dec=7"}
    q120["PrimitiveType__BlockEnd (120)<br/>BlockEnd<br/>"]

    q14 --> q119
    q113 -->|"tok(&quot;string&quot;)"| q114
    q114 --> q120
    q115 -->|"tok(&quot;bool&quot;)"| q116
    q116 --> q120
    q117 -->|"tok(&quot;composite&quot;)"| q118
    q118 --> q120
    q119 --> q113
    q119 --> q115
    q119 --> q117
    q120 --> q15
```

## ParserRule

```mermaid
flowchart TD
    q16(["ParserRule__Start (16)<br/>RuleStart"])
    q17(["ParserRule__Stop (17)<br/>RuleStop"])
    q121["ParserRule_Entry_entry (121)<br/>Basic<br/>"]
    q122["ParserRule__Basic_0 (122)<br/>Basic<br/>"]
    q123{"ParserRule__Basic_1 (123)<br/>Basic<br/><br/>dec=8"}
    q124["ParserRule_Name_ID (124)<br/>Basic<br/>"]
    q125["ParserRule_returns (125)<br/>Basic<br/>"]
    q126["ParserRule_ReturnType_ID (126)<br/>Basic<br/>"]
    q127["ParserRule__Basic_2 (127)<br/>Basic<br/>"]
    q128{"ParserRule__Basic_3 (128)<br/>Basic<br/><br/>dec=9"}
    q129["ParserRule_Colon (129)<br/>Basic<br/>"]
    q130["ParserRule__Basic_4 (130)<br/>Basic<br/>"]
    q131["ParserRule_Semicolon (131)<br/>Basic<br/>"]
    q132["ParserRule__Basic_5 (132)<br/>Basic<br/>"]
    q133{"ParserRule__Basic_6 (133)<br/>Basic<br/><br/>dec=10"}

    q16 --> q123
    q121 -->|"tok(&quot;entry&quot;)"| q122
    q122 --> q124
    q123 --> q121
    q123 --> q122
    q124 -->|"tok(ID)"| q128
    q125 -->|"tok(&quot;returns&quot;)"| q126
    q126 -->|"tok(ID)"| q127
    q127 --> q129
    q128 --> q125
    q128 --> q127
    q129 -->|"tok(&quot;:&quot;)"| q130
    q130 -.->|"[Alternatives]"| q133
    q131 -->|"tok(&quot;;&quot;)"| q132
    q132 --> q17
    q133 --> q131
    q133 --> q132
```

## Token

```mermaid
flowchart TD
    q18(["Token__Start (18)<br/>RuleStart"])
    q19(["Token__Stop (19)<br/>RuleStop"])
    q134["Token_Type_hidden (134)<br/>Basic<br/>"]
    q135["Token__Basic_0 (135)<br/>Basic<br/>"]
    q136["Token_Type_comment (136)<br/>Basic<br/>"]
    q137["Token__Basic_1 (137)<br/>Basic<br/>"]
    q138{"Token__Basic_2 (138)<br/>Basic<br/><br/>dec=11"}
    q139["Token__BlockEnd (139)<br/>BlockEnd<br/>"]
    q140{"Token__Basic_3 (140)<br/>Basic<br/><br/>dec=12"}
    q141["Token_token (141)<br/>Basic<br/>"]
    q142["Token_Name_ID (142)<br/>Basic<br/>"]
    q143["Token_Colon (143)<br/>Basic<br/>"]
    q144["Token_Regexp_RegexLiteral (144)<br/>Basic<br/>"]
    q145["Token_Semicolon (145)<br/>Basic<br/>"]
    q146["Token__Basic_4 (146)<br/>Basic<br/>"]
    q147{"Token__Basic_5 (147)<br/>Basic<br/><br/>dec=13"}

    q18 --> q140
    q134 -->|"tok(&quot;hidden&quot;)"| q135
    q135 --> q139
    q136 -->|"tok(&quot;comment&quot;)"| q137
    q137 --> q139
    q138 --> q134
    q138 --> q136
    q139 --> q141
    q140 --> q138
    q140 --> q139
    q141 -->|"tok(&quot;token&quot;)"| q142
    q142 -->|"tok(ID)"| q143
    q143 -->|"tok(&quot;:&quot;)"| q144
    q144 -->|"tok(RegexLiteral)"| q147
    q145 -->|"tok(&quot;;&quot;)"| q146
    q146 --> q19
    q147 --> q145
    q147 --> q146
```

## TokenGroup

```mermaid
flowchart TD
    q20(["TokenGroup__Start (20)<br/>RuleStart"])
    q21(["TokenGroup__Stop (21)<br/>RuleStop"])
    q148["TokenGroup_token (148)<br/>Basic<br/>"]
    q149["TokenGroup_group (149)<br/>Basic<br/>"]
    q150["TokenGroup_Name_ID (150)<br/>Basic<br/>"]
    q151["TokenGroup_LeftBrace (151)<br/>Basic<br/>"]
    q152["TokenGroup_TokenRefs_ID (152)<br/>Basic<br/>"]
    q153["TokenGroup__Basic_0 (153)<br/>Basic<br/>"]
    q154["TokenGroup_keywords (154)<br/>Basic<br/>"]
    q155["TokenGroup_Regexps_RegexLiteral (155)<br/>Basic<br/>"]
    q156["TokenGroup__Basic_1 (156)<br/>Basic<br/>"]
    q157["TokenGroup__Basic_2 (157)<br/>Basic<br/>"]
    q158["TokenGroup__Basic_3 (158)<br/>Basic<br/>"]
    q159{"TokenGroup__Basic_4 (159)<br/>Basic<br/><br/>dec=14"}
    q160["TokenGroup__BlockEnd (160)<br/>BlockEnd<br/>"]
    q161{"TokenGroup__LoopEntry (161)<br/>LoopEntry<br/><br/>dec=15"}
    q162["TokenGroup__LoopEnd (162)<br/>LoopEnd<br/>"]
    q163["TokenGroup__LoopBack (163)<br/>LoopBack<br/>"]
    q164["TokenGroup_RightBrace (164)<br/>Basic<br/>"]
    q165["TokenGroup__Basic_5 (165)<br/>Basic<br/>"]

    q20 --> q148
    q148 -->|"tok(&quot;token&quot;)"| q149
    q149 -->|"tok(&quot;group&quot;)"| q150
    q150 -->|"tok(ID)"| q151
    q151 -->|"tok(&quot;{&quot;)"| q161
    q152 -->|"tok(ID)"| q153
    q153 --> q160
    q154 -->|"tok(&quot;keywords&quot;)"| q155
    q155 -->|"tok(RegexLiteral)"| q156
    q156 --> q160
    q157 -.->|"[Keyword]"| q158
    q158 --> q160
    q159 --> q152
    q159 --> q154
    q159 --> q157
    q160 --> q163
    q161 --> q159
    q161 --> q162
    q162 --> q164
    q163 --> q161
    q164 -->|"tok(&quot;}&quot;)"| q165
    q165 --> q21
```

## Alternatives

```mermaid
flowchart TD
    q22(["Alternatives__Start (22)<br/>RuleStart"])
    q23(["Alternatives__Stop (23)<br/>RuleStop"])
    q166["Alternatives__Basic_0 (166)<br/>Basic<br/>"]
    q167["Alternatives_Pipe (167)<br/>Basic<br/>"]
    q168["Alternatives__Basic_1 (168)<br/>Basic<br/>"]
    q169["Alternatives__Basic_2 (169)<br/>Basic<br/>"]
    q170{"Alternatives__LoopBack (170)<br/>LoopBack<br/><br/>dec=16"}
    q171["Alternatives__LoopEnd (171)<br/>LoopEnd<br/>"]
    q172{"Alternatives__Basic_3 (172)<br/>Basic<br/><br/>dec=17"}

    q22 --> q166
    q166 -.->|"[Group]"| q172
    q167 -->|"tok(&quot;|&quot;)"| q168
    q168 -.->|"[Group]"| q169
    q169 --> q170
    q170 --> q167
    q170 --> q171
    q171 --> q23
    q172 --> q167
    q172 --> q171
```

## Group

```mermaid
flowchart TD
    q24(["Group__Start (24)<br/>RuleStart"])
    q25(["Group__Stop (25)<br/>RuleStop"])
    q173["Group__Basic_0 (173)<br/>Basic<br/>"]
    q174["Group__Basic_1 (174)<br/>Basic<br/>"]
    q175["Group__Basic_2 (175)<br/>Basic<br/>"]
    q176{"Group__LoopBack (176)<br/>LoopBack<br/><br/>dec=18"}
    q177["Group__LoopEnd (177)<br/>LoopEnd<br/>"]
    q178{"Group__Basic_3 (178)<br/>Basic<br/><br/>dec=19"}

    q24 --> q173
    q173 -.->|"[Element]"| q178
    q174 -.->|"[Element]"| q175
    q175 --> q176
    q176 --> q174
    q176 --> q177
    q177 --> q25
    q178 --> q174
    q178 --> q177
```

## Element

```mermaid
flowchart TD
    q26(["Element__Start (26)<br/>RuleStart"])
    q27(["Element__Stop (27)<br/>RuleStop"])
    q179["Element__Basic_0 (179)<br/>Basic<br/>"]
    q180["Element__Basic_1 (180)<br/>Basic<br/>"]
    q181["Element__Basic_2 (181)<br/>Basic<br/>"]
    q182["Element__Basic_3 (182)<br/>Basic<br/>"]
    q183["Element__Basic_4 (183)<br/>Basic<br/>"]
    q184["Element__Basic_5 (184)<br/>Basic<br/>"]
    q185["Element__Basic_6 (185)<br/>Basic<br/>"]
    q186["Element__Basic_7 (186)<br/>Basic<br/>"]
    q187["Element_LeftParen (187)<br/>Basic<br/>"]
    q188["Element__Basic_8 (188)<br/>Basic<br/>"]
    q189["Element_RightParen (189)<br/>Basic<br/>"]
    q190["Element__Basic_9 (190)<br/>Basic<br/>"]
    q191{"Element__Basic_10 (191)<br/>Basic<br/><br/>dec=20"}
    q192["Element__BlockEnd_0 (192)<br/>BlockEnd<br/>"]
    q193["Element_Cardinality_Asterisk (193)<br/>Basic<br/>"]
    q194["Element__Basic_11 (194)<br/>Basic<br/>"]
    q195["Element_Cardinality_Plus (195)<br/>Basic<br/>"]
    q196["Element__Basic_12 (196)<br/>Basic<br/>"]
    q197["Element_Cardinality_Question (197)<br/>Basic<br/>"]
    q198["Element__Basic_13 (198)<br/>Basic<br/>"]
    q199{"Element__Basic_14 (199)<br/>Basic<br/><br/>dec=21"}
    q200["Element__BlockEnd_1 (200)<br/>BlockEnd<br/>"]
    q201{"Element__Basic_15 (201)<br/>Basic<br/><br/>dec=22"}

    q26 --> q191
    q179 -.->|"[Keyword]"| q180
    q180 --> q192
    q181 -.->|"[Assignment]"| q182
    q182 --> q192
    q183 -.->|"[RuleCall]"| q184
    q184 --> q192
    q185 -.->|"[Action]"| q186
    q186 --> q192
    q187 -->|"tok(&quot;(&quot;)"| q188
    q188 -.->|"[Alternatives]"| q189
    q189 -->|"tok(&quot;)&quot;)"| q190
    q190 --> q192
    q191 --> q179
    q191 --> q181
    q191 --> q183
    q191 --> q185
    q191 --> q187
    q192 --> q201
    q193 -->|"tok(&quot;*&quot;)"| q194
    q194 --> q200
    q195 -->|"tok(&quot;+&quot;)"| q196
    q196 --> q200
    q197 -->|"tok(&quot;?&quot;)"| q198
    q198 --> q200
    q199 --> q193
    q199 --> q195
    q199 --> q197
    q200 --> q27
    q201 --> q199
    q201 --> q200
```

## Keyword

```mermaid
flowchart TD
    q28(["Keyword__Start (28)<br/>RuleStart"])
    q29(["Keyword__Stop (29)<br/>RuleStop"])
    q202["Keyword_Value_StringLiteral (202)<br/>Basic<br/>"]
    q203["Keyword__Basic (203)<br/>Basic<br/>"]

    q28 --> q202
    q202 -->|"tok(StringLiteral)"| q203
    q203 --> q29
```

## Assignment

```mermaid
flowchart TD
    q30(["Assignment__Start (30)<br/>RuleStart"])
    q31(["Assignment__Stop (31)<br/>RuleStop"])
    q204["Assignment_Property_ID (204)<br/>Basic<br/>"]
    q205["Assignment_Operator_PlusEquals (205)<br/>Basic<br/>"]
    q206["Assignment__Basic_0 (206)<br/>Basic<br/>"]
    q207["Assignment_Operator_Equals (207)<br/>Basic<br/>"]
    q208["Assignment__Basic_1 (208)<br/>Basic<br/>"]
    q209["Assignment_Operator_QuestionEquals (209)<br/>Basic<br/>"]
    q210["Assignment__Basic_2 (210)<br/>Basic<br/>"]
    q211{"Assignment__Basic_3 (211)<br/>Basic<br/><br/>dec=23"}
    q212["Assignment__BlockEnd (212)<br/>BlockEnd<br/>"]
    q213["Assignment__Basic_4 (213)<br/>Basic<br/>"]
    q214["Assignment__Basic_5 (214)<br/>Basic<br/>"]

    q30 --> q204
    q204 -->|"tok(ID)"| q211
    q205 -->|"tok(&quot;+=&quot;)"| q206
    q206 --> q212
    q207 -->|"tok(&quot;=&quot;)"| q208
    q208 --> q212
    q209 -->|"tok(&quot;?=&quot;)"| q210
    q210 --> q212
    q211 --> q205
    q211 --> q207
    q211 --> q209
    q212 --> q213
    q213 -.->|"[Assignable]"| q214
    q214 --> q31
```

## Assignable

```mermaid
flowchart TD
    q32(["Assignable__Start (32)<br/>RuleStart"])
    q33(["Assignable__Stop (33)<br/>RuleStop"])
    q215["Assignable__Basic_0 (215)<br/>Basic<br/>"]
    q216["Assignable__Basic_1 (216)<br/>Basic<br/>"]
    q217["Assignable__Basic_2 (217)<br/>Basic<br/>"]
    q218["Assignable__Basic_3 (218)<br/>Basic<br/>"]
    q219["Assignable__Basic_4 (219)<br/>Basic<br/>"]
    q220["Assignable__Basic_5 (220)<br/>Basic<br/>"]
    q221["Assignable_LeftParen (221)<br/>Basic<br/>"]
    q222["Assignable__Basic_6 (222)<br/>Basic<br/>"]
    q223["Assignable_RightParen (223)<br/>Basic<br/>"]
    q224["Assignable__Basic_7 (224)<br/>Basic<br/>"]
    q225{"Assignable__Basic_8 (225)<br/>Basic<br/><br/>dec=24"}
    q226["Assignable__BlockEnd (226)<br/>BlockEnd<br/>"]

    q32 --> q225
    q215 -.->|"[Keyword]"| q216
    q216 --> q226
    q217 -.->|"[RuleCall]"| q218
    q218 --> q226
    q219 -.->|"[CrossRef]"| q220
    q220 --> q226
    q221 -->|"tok(&quot;(&quot;)"| q222
    q222 -.->|"[AssignableAlternatives]"| q223
    q223 -->|"tok(&quot;)&quot;)"| q224
    q224 --> q226
    q225 --> q215
    q225 --> q217
    q225 --> q219
    q225 --> q221
    q226 --> q33
```

## AssignableWithoutAlts

```mermaid
flowchart TD
    q34(["AssignableWithoutAlts__Start (34)<br/>RuleStart"])
    q35(["AssignableWithoutAlts__Stop (35)<br/>RuleStop"])
    q227["AssignableWithoutAlts__Basic_0 (227)<br/>Basic<br/>"]
    q228["AssignableWithoutAlts__Basic_1 (228)<br/>Basic<br/>"]
    q229["AssignableWithoutAlts__Basic_2 (229)<br/>Basic<br/>"]
    q230["AssignableWithoutAlts__Basic_3 (230)<br/>Basic<br/>"]
    q231["AssignableWithoutAlts__Basic_4 (231)<br/>Basic<br/>"]
    q232["AssignableWithoutAlts__Basic_5 (232)<br/>Basic<br/>"]
    q233{"AssignableWithoutAlts__Basic_6 (233)<br/>Basic<br/><br/>dec=25"}
    q234["AssignableWithoutAlts__BlockEnd (234)<br/>BlockEnd<br/>"]

    q34 --> q233
    q227 -.->|"[Keyword]"| q228
    q228 --> q234
    q229 -.->|"[RuleCall]"| q230
    q230 --> q234
    q231 -.->|"[CrossRef]"| q232
    q232 --> q234
    q233 --> q227
    q233 --> q229
    q233 --> q231
    q234 --> q35
```

## AssignableAlternatives

```mermaid
flowchart TD
    q36(["AssignableAlternatives__Start (36)<br/>RuleStart"])
    q37(["AssignableAlternatives__Stop (37)<br/>RuleStop"])
    q235["AssignableAlternatives__Basic_0 (235)<br/>Basic<br/>"]
    q236["AssignableAlternatives_Pipe (236)<br/>Basic<br/>"]
    q237["AssignableAlternatives__Basic_1 (237)<br/>Basic<br/>"]
    q238["AssignableAlternatives__Basic_2 (238)<br/>Basic<br/>"]
    q239{"AssignableAlternatives__LoopBack (239)<br/>LoopBack<br/><br/>dec=26"}
    q240["AssignableAlternatives__LoopEnd (240)<br/>LoopEnd<br/>"]
    q241{"AssignableAlternatives__Basic_3 (241)<br/>Basic<br/><br/>dec=27"}

    q36 --> q235
    q235 -.->|"[AssignableWithoutAlts]"| q241
    q236 -->|"tok(&quot;|&quot;)"| q237
    q237 -.->|"[AssignableWithoutAlts]"| q238
    q238 --> q239
    q239 --> q236
    q239 --> q240
    q240 --> q37
    q241 --> q236
    q241 --> q240
```

## CrossRef

```mermaid
flowchart TD
    q38(["CrossRef__Start (38)<br/>RuleStart"])
    q39(["CrossRef__Stop (39)<br/>RuleStop"])
    q242["CrossRef_LeftBracket (242)<br/>Basic<br/>"]
    q243["CrossRef_Type_ID (243)<br/>Basic<br/>"]
    q244["CrossRef_Colon (244)<br/>Basic<br/>"]
    q245["CrossRef__Basic_0 (245)<br/>Basic<br/>"]
    q246["CrossRef__Basic_1 (246)<br/>Basic<br/>"]
    q247{"CrossRef__Basic_2 (247)<br/>Basic<br/><br/>dec=28"}
    q248["CrossRef_RightBracket (248)<br/>Basic<br/>"]
    q249["CrossRef__Basic_3 (249)<br/>Basic<br/>"]

    q38 --> q242
    q242 -->|"tok(&quot;[&quot;)"| q243
    q243 -->|"tok(ID)"| q247
    q244 -->|"tok(&quot;:&quot;)"| q245
    q245 -.->|"[RuleCall]"| q246
    q246 --> q248
    q247 --> q244
    q247 --> q246
    q248 -->|"tok(&quot;]&quot;)"| q249
    q249 --> q39
```

## RuleCall

```mermaid
flowchart TD
    q40(["RuleCall__Start (40)<br/>RuleStart"])
    q41(["RuleCall__Stop (41)<br/>RuleStop"])
    q250["RuleCall_Rule_ID (250)<br/>Basic<br/>"]
    q251["RuleCall__Basic (251)<br/>Basic<br/>"]

    q40 --> q250
    q250 -->|"tok(ID)"| q251
    q251 --> q41
```

## Action

```mermaid
flowchart TD
    q42(["Action__Start (42)<br/>RuleStart"])
    q43(["Action__Stop (43)<br/>RuleStop"])
    q252["Action_LeftBrace (252)<br/>Basic<br/>"]
    q253["Action_Type_ID (253)<br/>Basic<br/>"]
    q254["Action_Dot (254)<br/>Basic<br/>"]
    q255["Action_Property_ID (255)<br/>Basic<br/>"]
    q256["Action_Operator_PlusEquals (256)<br/>Basic<br/>"]
    q257["Action__Basic_0 (257)<br/>Basic<br/>"]
    q258["Action_Operator_Equals (258)<br/>Basic<br/>"]
    q259["Action__Basic_1 (259)<br/>Basic<br/>"]
    q260{"Action__Basic_2 (260)<br/>Basic<br/><br/>dec=29"}
    q261["Action__BlockEnd (261)<br/>BlockEnd<br/>"]
    q262["Action_current (262)<br/>Basic<br/>"]
    q263["Action__Basic_3 (263)<br/>Basic<br/>"]
    q264{"Action__Basic_4 (264)<br/>Basic<br/><br/>dec=30"}
    q265["Action_RightBrace (265)<br/>Basic<br/>"]
    q266["Action__Basic_5 (266)<br/>Basic<br/>"]

    q42 --> q252
    q252 -->|"tok(&quot;{&quot;)"| q253
    q253 -->|"tok(ID)"| q264
    q254 -->|"tok(&quot;.&quot;)"| q255
    q255 -->|"tok(ID)"| q260
    q256 -->|"tok(&quot;+=&quot;)"| q257
    q257 --> q261
    q258 -->|"tok(&quot;=&quot;)"| q259
    q259 --> q261
    q260 --> q256
    q260 --> q258
    q261 --> q262
    q262 -->|"tok(&quot;current&quot;)"| q263
    q263 --> q265
    q264 --> q254
    q264 --> q263
    q265 -->|"tok(&quot;}&quot;)"| q266
    q266 --> q43
```

## CompositeRule

```mermaid
flowchart TD
    q44(["CompositeRule__Start (44)<br/>RuleStart"])
    q45(["CompositeRule__Stop (45)<br/>RuleStop"])
    q267["CompositeRule_composite (267)<br/>Basic<br/>"]
    q268["CompositeRule_Name_ID (268)<br/>Basic<br/>"]
    q269["CompositeRule_Colon (269)<br/>Basic<br/>"]
    q270["CompositeRule__Basic_0 (270)<br/>Basic<br/>"]
    q271["CompositeRule_Semicolon (271)<br/>Basic<br/>"]
    q272["CompositeRule__Basic_1 (272)<br/>Basic<br/>"]
    q273{"CompositeRule__Basic_2 (273)<br/>Basic<br/><br/>dec=31"}

    q44 --> q267
    q267 -->|"tok(&quot;composite&quot;)"| q268
    q268 -->|"tok(ID)"| q269
    q269 -->|"tok(&quot;:&quot;)"| q270
    q270 -.->|"[CompositeAlternatives]"| q273
    q271 -->|"tok(&quot;;&quot;)"| q272
    q272 --> q45
    q273 --> q271
    q273 --> q272
```

## CompositeAlternatives

```mermaid
flowchart TD
    q46(["CompositeAlternatives__Start (46)<br/>RuleStart"])
    q47(["CompositeAlternatives__Stop (47)<br/>RuleStop"])
    q274["CompositeAlternatives__Basic_0 (274)<br/>Basic<br/>"]
    q275["CompositeAlternatives_Pipe (275)<br/>Basic<br/>"]
    q276["CompositeAlternatives__Basic_1 (276)<br/>Basic<br/>"]
    q277["CompositeAlternatives__Basic_2 (277)<br/>Basic<br/>"]
    q278{"CompositeAlternatives__LoopBack (278)<br/>LoopBack<br/><br/>dec=32"}
    q279["CompositeAlternatives__LoopEnd (279)<br/>LoopEnd<br/>"]
    q280{"CompositeAlternatives__Basic_3 (280)<br/>Basic<br/><br/>dec=33"}

    q46 --> q274
    q274 -.->|"[CompositeGroup]"| q280
    q275 -->|"tok(&quot;|&quot;)"| q276
    q276 -.->|"[CompositeGroup]"| q277
    q277 --> q278
    q278 --> q275
    q278 --> q279
    q279 --> q47
    q280 --> q275
    q280 --> q279
```

## CompositeGroup

```mermaid
flowchart TD
    q48(["CompositeGroup__Start (48)<br/>RuleStart"])
    q49(["CompositeGroup__Stop (49)<br/>RuleStop"])
    q281["CompositeGroup__Basic_0 (281)<br/>Basic<br/>"]
    q282["CompositeGroup__Basic_1 (282)<br/>Basic<br/>"]
    q283["CompositeGroup__Basic_2 (283)<br/>Basic<br/>"]
    q284{"CompositeGroup__LoopBack (284)<br/>LoopBack<br/><br/>dec=34"}
    q285["CompositeGroup__LoopEnd (285)<br/>LoopEnd<br/>"]
    q286{"CompositeGroup__Basic_3 (286)<br/>Basic<br/><br/>dec=35"}

    q48 --> q281
    q281 -.->|"[CompositeElement]"| q286
    q282 -.->|"[CompositeElement]"| q283
    q283 --> q284
    q284 --> q282
    q284 --> q285
    q285 --> q49
    q286 --> q282
    q286 --> q285
```

## CompositeElement

```mermaid
flowchart TD
    q50(["CompositeElement__Start (50)<br/>RuleStart"])
    q51(["CompositeElement__Stop (51)<br/>RuleStop"])
    q287["CompositeElement__Basic_0 (287)<br/>Basic<br/>"]
    q288["CompositeElement__Basic_1 (288)<br/>Basic<br/>"]
    q289["CompositeElement__Basic_2 (289)<br/>Basic<br/>"]
    q290["CompositeElement__Basic_3 (290)<br/>Basic<br/>"]
    q291["CompositeElement_LeftParen (291)<br/>Basic<br/>"]
    q292["CompositeElement__Basic_4 (292)<br/>Basic<br/>"]
    q293["CompositeElement_RightParen (293)<br/>Basic<br/>"]
    q294["CompositeElement__Basic_5 (294)<br/>Basic<br/>"]
    q295{"CompositeElement__Basic_6 (295)<br/>Basic<br/><br/>dec=36"}
    q296["CompositeElement__BlockEnd_0 (296)<br/>BlockEnd<br/>"]
    q297["CompositeElement_Cardinality_Asterisk (297)<br/>Basic<br/>"]
    q298["CompositeElement__Basic_7 (298)<br/>Basic<br/>"]
    q299["CompositeElement_Cardinality_Plus (299)<br/>Basic<br/>"]
    q300["CompositeElement__Basic_8 (300)<br/>Basic<br/>"]
    q301["CompositeElement_Cardinality_Question (301)<br/>Basic<br/>"]
    q302["CompositeElement__Basic_9 (302)<br/>Basic<br/>"]
    q303{"CompositeElement__Basic_10 (303)<br/>Basic<br/><br/>dec=37"}
    q304["CompositeElement__BlockEnd_1 (304)<br/>BlockEnd<br/>"]
    q305{"CompositeElement__Basic_11 (305)<br/>Basic<br/><br/>dec=38"}

    q50 --> q295
    q287 -.->|"[Keyword]"| q288
    q288 --> q296
    q289 -.->|"[RuleCall]"| q290
    q290 --> q296
    q291 -->|"tok(&quot;(&quot;)"| q292
    q292 -.->|"[CompositeAlternatives]"| q293
    q293 -->|"tok(&quot;)&quot;)"| q294
    q294 --> q296
    q295 --> q287
    q295 --> q289
    q295 --> q291
    q296 --> q305
    q297 -->|"tok(&quot;*&quot;)"| q298
    q298 --> q304
    q299 -->|"tok(&quot;+&quot;)"| q300
    q300 --> q304
    q301 -->|"tok(&quot;?&quot;)"| q302
    q302 --> q304
    q303 --> q297
    q303 --> q299
    q303 --> q301
    q304 --> q51
    q305 --> q303
    q305 --> q304
```

