# Runtime ATN for completion

## Root

```mermaid
flowchart TD
    q0(["Root__Start (0)<br/>RuleStart"])
    q1(["Root__Stop (1)<br/>RuleStop"])
    q46["Root__Basic_0 (46)<br/>Basic<br/>"]
    q47["Root__Basic_1 (47)<br/>Basic<br/>"]
    q48["Root__Basic_2 (48)<br/>Basic<br/>"]
    q49["Root__Basic_3 (49)<br/>Basic<br/>"]
    q50["Root__Basic_4 (50)<br/>Basic<br/>"]
    q51["Root__Basic_5 (51)<br/>Basic<br/>"]
    q52["Root__Basic_6 (52)<br/>Basic<br/>"]
    q53["Root__Basic_7 (53)<br/>Basic<br/>"]
    q54["Root__Basic_8 (54)<br/>Basic<br/>"]
    q55["Root__Basic_9 (55)<br/>Basic<br/>"]
    q56["Root__Basic_10 (56)<br/>Basic<br/>"]
    q57["Root__Basic_11 (57)<br/>Basic<br/>"]
    q58["Root__Basic_12 (58)<br/>Basic<br/>"]
    q59["Root__Basic_13 (59)<br/>Basic<br/>"]
    q60["Root__Basic_14 (60)<br/>Basic<br/>"]
    q61["Root__Basic_15 (61)<br/>Basic<br/>"]
    q62["Root__Basic_16 (62)<br/>Basic<br/>"]
    q63["Root__Basic_17 (63)<br/>Basic<br/>"]
    q64["Root__Basic_18 (64)<br/>Basic<br/>"]
    q65["Root__Basic_19 (65)<br/>Basic<br/>"]
    q66["Root__Basic_20 (66)<br/>Basic<br/>"]
    q67["Root__Basic_21 (67)<br/>Basic<br/>"]
    q68["Root__Basic_22 (68)<br/>Basic<br/>"]
    q69["Root__Basic_23 (69)<br/>Basic<br/>"]
    q70["Root__Basic_24 (70)<br/>Basic<br/>"]
    q71["Root__Basic_25 (71)<br/>Basic<br/>"]
    q72["Root__Basic_26 (72)<br/>Basic<br/>"]
    q73["Root__Basic_27 (73)<br/>Basic<br/>"]
    q74["Root__Basic_28 (74)<br/>Basic<br/>"]
    q75["Root__Basic_29 (75)<br/>Basic<br/>"]
    q76["Root__Basic_30 (76)<br/>Basic<br/>"]
    q77["Root__Basic_31 (77)<br/>Basic<br/>"]
    q78{"Root__Basic_32 (78)<br/>Basic<br/><br/>dec=0"}
    q79["Root__BlockEnd (79)<br/>BlockEnd<br/>"]
    q80{"Root__LoopEntry (80)<br/>LoopEntry<br/><br/>dec=1"}
    q81["Root__LoopEnd (81)<br/>LoopEnd<br/>"]
    q82["Root__LoopBack (82)<br/>LoopBack<br/>"]

    q0 --> q80
    q46 -.->|"[Declare]"| q47
    q47 --> q79
    q48 -.->|"[A]"| q49
    q49 --> q79
    q50 -.->|"[B]"| q51
    q51 --> q79
    q52 -.->|"[C]"| q53
    q53 --> q79
    q54 -.->|"[D]"| q55
    q55 --> q79
    q56 -.->|"[E]"| q57
    q57 --> q79
    q58 -.->|"[F]"| q59
    q59 --> q79
    q60 -.->|"[G]"| q61
    q61 --> q79
    q62 -.->|"[H]"| q63
    q63 --> q79
    q64 -.->|"[I]"| q65
    q65 --> q79
    q66 -.->|"[J]"| q67
    q67 --> q79
    q68 -.->|"[K]"| q69
    q69 --> q79
    q70 -.->|"[L]"| q71
    q71 --> q79
    q72 -.->|"[M]"| q73
    q73 --> q79
    q74 -.->|"[N]"| q75
    q75 --> q79
    q76 -.->|"[O]"| q77
    q77 --> q79
    q78 --> q46
    q78 --> q48
    q78 --> q50
    q78 --> q52
    q78 --> q54
    q78 --> q56
    q78 --> q58
    q78 --> q60
    q78 --> q62
    q78 --> q64
    q78 --> q66
    q78 --> q68
    q78 --> q70
    q78 --> q72
    q78 --> q74
    q78 --> q76
    q79 --> q82
    q80 --> q78
    q80 --> q81
    q81 --> q1
    q82 --> q80
```

## Declare

```mermaid
flowchart TD
    q2(["Declare__Start (2)<br/>RuleStart"])
    q3(["Declare__Stop (3)<br/>RuleStop"])
    q83["Declare_DECLARE (83)<br/>Basic<br/>"]
    q84["Declare__Basic_0 (84)<br/>Basic<br/>"]
    q85["Declare_LBRACE (85)<br/>Basic<br/>"]
    q86["Declare__Basic_1 (86)<br/>Basic<br/>"]
    q87["Declare__Basic_2 (87)<br/>Basic<br/>"]
    q88{"Declare__LoopEntry (88)<br/>LoopEntry<br/><br/>dec=2"}
    q89["Declare__LoopEnd (89)<br/>LoopEnd<br/>"]
    q90["Declare__LoopBack (90)<br/>LoopBack<br/>"]
    q91["Declare_RBRACE (91)<br/>Basic<br/>"]
    q92["Declare__Basic_3 (92)<br/>Basic<br/>"]
    q93{"Declare__Basic_4 (93)<br/>Basic<br/><br/>dec=3"}

    q2 --> q83
    q83 -->|"tok(Token_DECLARE)"| q84
    q84 -.->|"[FQN]"| q93
    q85 -->|"tok(Token_LBRACE)"| q88
    q86 -.->|"[Declare]"| q87
    q87 --> q90
    q88 --> q86
    q88 --> q89
    q89 --> q91
    q90 --> q88
    q91 -->|"tok(Token_RBRACE)"| q92
    q92 --> q3
    q93 --> q85
    q93 --> q92
```

## A

```mermaid
flowchart TD
    q4(["A__Start (4)<br/>RuleStart"])
    q5(["A__Stop (5)<br/>RuleStop"])
    q94["A_a (94)<br/>Basic<br/>"]
    q95["A_FIRST (95)<br/>Basic<br/>"]
    q96["A__Basic (96)<br/>Basic<br/>"]

    q4 --> q94
    q94 -->|"tok(Keyword_a)"| q95
    q95 -->|"tok(Token_FIRST)"| q96
    q96 --> q5
```

## B

```mermaid
flowchart TD
    q6(["B__Start (6)<br/>RuleStart"])
    q7(["B__Stop (7)<br/>RuleStop"])
    q97["B_b (97)<br/>Basic<br/>"]
    q98["B_FIRST (98)<br/>Basic<br/>"]
    q99["B__Basic_0 (99)<br/>Basic<br/>"]
    q100["B_SECOND (100)<br/>Basic<br/>"]
    q101["B__Basic_1 (101)<br/>Basic<br/>"]
    q102{"B__Basic_2 (102)<br/>Basic<br/><br/>dec=4"}
    q103["B__BlockEnd (103)<br/>BlockEnd<br/>"]

    q6 --> q97
    q97 -->|"tok(Keyword_b)"| q102
    q98 -->|"tok(Token_FIRST)"| q99
    q99 --> q103
    q100 -->|"tok(Token_SECOND)"| q101
    q101 --> q103
    q102 --> q98
    q102 --> q100
    q103 --> q7
```

## C

```mermaid
flowchart TD
    q8(["C__Start (8)<br/>RuleStart"])
    q9(["C__Stop (9)<br/>RuleStop"])
    q104["C_c (104)<br/>Basic<br/>"]
    q105["C_COMMON_0 (105)<br/>Basic<br/>"]
    q106["C_FIRST (106)<br/>Basic<br/>"]
    q107["C__Basic_0 (107)<br/>Basic<br/>"]
    q108["C_COMMON_1 (108)<br/>Basic<br/>"]
    q109["C_SECOND (109)<br/>Basic<br/>"]
    q110["C__Basic_1 (110)<br/>Basic<br/>"]
    q111{"C__Basic_2 (111)<br/>Basic<br/><br/>dec=5"}
    q112["C__BlockEnd (112)<br/>BlockEnd<br/>"]

    q8 --> q104
    q104 -->|"tok(Keyword_c)"| q111
    q105 -->|"tok(Token_COMMON)"| q106
    q106 -->|"tok(Token_FIRST)"| q107
    q107 --> q112
    q108 -->|"tok(Token_COMMON)"| q109
    q109 -->|"tok(Token_SECOND)"| q110
    q110 --> q112
    q111 --> q105
    q111 --> q108
    q112 --> q9
```

## D

```mermaid
flowchart TD
    q10(["D__Start (10)<br/>RuleStart"])
    q11(["D__Stop (11)<br/>RuleStop"])
    q113["D_d (113)<br/>Basic<br/>"]
    q114["D__Basic_0 (114)<br/>Basic<br/>"]
    q115["D__Basic_1 (115)<br/>Basic<br/>"]
    q116["D__Basic_2 (116)<br/>Basic<br/>"]
    q117["D__Basic_3 (117)<br/>Basic<br/>"]
    q118{"D__Basic_4 (118)<br/>Basic<br/><br/>dec=6"}
    q119["D__BlockEnd (119)<br/>BlockEnd<br/>"]

    q10 --> q113
    q113 -->|"tok(Keyword_d)"| q118
    q114 -.->|"[DLong]"| q115
    q115 --> q119
    q116 -.->|"[DShort]"| q117
    q117 --> q119
    q118 --> q114
    q118 --> q116
    q119 --> q11
```

## E

```mermaid
flowchart TD
    q12(["E__Start (12)<br/>RuleStart"])
    q13(["E__Stop (13)<br/>RuleStop"])
    q120["E_e (120)<br/>Basic<br/>"]
    q121["E__Basic_0 (121)<br/>Basic<br/>"]
    q122["E__Basic_1 (122)<br/>Basic<br/>"]

    q12 --> q120
    q120 -->|"tok(Keyword_e)"| q121
    q121 -.->|"[FQN]"| q122
    q122 --> q13
```

## DLong

```mermaid
flowchart TD
    q14(["DLong__Start (14)<br/>RuleStart"])
    q15(["DLong__Stop (15)<br/>RuleStop"])
    q123["DLong_COMMON (123)<br/>Basic<br/>"]
    q124["DLong_THEN (124)<br/>Basic<br/>"]
    q125["DLong_LONG (125)<br/>Basic<br/>"]
    q126["DLong__Basic (126)<br/>Basic<br/>"]

    q14 --> q123
    q123 -->|"tok(Token_COMMON)"| q124
    q124 -->|"tok(Token_THEN)"| q125
    q125 -->|"tok(Token_LONG)"| q126
    q126 --> q15
```

## DShort

```mermaid
flowchart TD
    q16(["DShort__Start (16)<br/>RuleStart"])
    q17(["DShort__Stop (17)<br/>RuleStop"])
    q127["DShort_COMMON (127)<br/>Basic<br/>"]
    q128["DShort__Basic (128)<br/>Basic<br/>"]

    q16 --> q127
    q127 -->|"tok(Token_COMMON)"| q128
    q128 --> q17
```

## F

```mermaid
flowchart TD
    q18(["F__Start (18)<br/>RuleStart"])
    q19(["F__Stop (19)<br/>RuleStop"])
    q129["F_f (129)<br/>Basic<br/>"]
    q130["F__Basic_0 (130)<br/>Basic<br/>"]
    q131["F__Basic_1 (131)<br/>Basic<br/>"]
    q132{"F__LoopBack (132)<br/>LoopBack<br/><br/>dec=7"}
    q133["F__LoopEnd (133)<br/>LoopEnd<br/>"]

    q18 --> q129
    q129 -->|"tok(Keyword_f)"| q130
    q130 -.->|"[FItem]"| q131
    q131 --> q132
    q132 --> q130
    q132 --> q133
    q133 --> q19
```

## FItem

```mermaid
flowchart TD
    q20(["FItem__Start (20)<br/>RuleStart"])
    q21(["FItem__Stop (21)<br/>RuleStop"])
    q134["FItem__Basic_0 (134)<br/>Basic<br/>"]
    q135["FItem__Basic_1 (135)<br/>Basic<br/>"]

    q20 --> q134
    q134 -.->|"[FQN]"| q135
    q135 --> q21
```

## G

```mermaid
flowchart TD
    q22(["G__Start (22)<br/>RuleStart"])
    q23(["G__Stop (23)<br/>RuleStop"])
    q136["G_g (136)<br/>Basic<br/>"]
    q137["G_Ref_ID (137)<br/>Basic<br/>"]
    q138["G__Basic (138)<br/>Basic<br/>"]

    q22 --> q136
    q136 -->|"tok(Keyword_g)"| q137
    q137 -->|"tok(Token_ID)"| q138
    q138 --> q23
```

## H

```mermaid
flowchart TD
    q24(["H__Start (24)<br/>RuleStart"])
    q25(["H__Stop (25)<br/>RuleStop"])
    q139["H_h (139)<br/>Basic<br/>"]
    q140["H__Basic_0 (140)<br/>Basic<br/>"]
    q141["H__Basic_1 (141)<br/>Basic<br/>"]

    q24 --> q139
    q139 -->|"tok(Keyword_h)"| q140
    q140 -.->|"[MemberCall]"| q141
    q141 --> q25
```

## I

```mermaid
flowchart TD
    q26(["I__Start (26)<br/>RuleStart"])
    q27(["I__Stop (27)<br/>RuleStop"])
    q142["I_i (142)<br/>Basic<br/>"]
    q143["I__Basic_0 (143)<br/>Basic<br/>"]
    q144["I__Basic_1 (144)<br/>Basic<br/>"]

    q26 --> q142
    q142 -->|"tok(Keyword_i)"| q143
    q143 -.->|"[MemberCallNoDot]"| q144
    q144 --> q27
```

## MemberCall

```mermaid
flowchart TD
    q28(["MemberCall__Start (28)<br/>RuleStart"])
    q29(["MemberCall__Stop (29)<br/>RuleStop"])
    q145["MemberCall_Ref_ID_0 (145)<br/>Basic<br/>"]
    q146["MemberCall_DOT (146)<br/>Basic<br/>"]
    q147["MemberCall_Ref_ID_1 (147)<br/>Basic<br/>"]
    q148["MemberCall__Basic (148)<br/>Basic<br/>"]
    q149{"MemberCall__LoopEntry (149)<br/>LoopEntry<br/><br/>dec=8"}
    q150["MemberCall__LoopEnd (150)<br/>LoopEnd<br/>"]
    q151["MemberCall__LoopBack (151)<br/>LoopBack<br/>"]

    q28 --> q145
    q145 -->|"tok(Token_ID)"| q149
    q146 -->|"tok(Token_DOT)"| q147
    q147 -->|"tok(Token_ID)"| q148
    q148 --> q151
    q149 --> q146
    q149 --> q150
    q150 --> q29
    q151 --> q149
```

## MemberCallNoDot

```mermaid
flowchart TD
    q30(["MemberCallNoDot__Start (30)<br/>RuleStart"])
    q31(["MemberCallNoDot__Stop (31)<br/>RuleStop"])
    q152["MemberCallNoDot_Ref_ID_0 (152)<br/>Basic<br/>"]
    q153["MemberCallNoDot_Ref_ID_1 (153)<br/>Basic<br/>"]
    q154["MemberCallNoDot__Basic (154)<br/>Basic<br/>"]
    q155{"MemberCallNoDot__LoopEntry (155)<br/>LoopEntry<br/><br/>dec=9"}
    q156["MemberCallNoDot__LoopEnd (156)<br/>LoopEnd<br/>"]
    q157["MemberCallNoDot__LoopBack (157)<br/>LoopBack<br/>"]

    q30 --> q152
    q152 -->|"tok(Token_ID)"| q155
    q153 -->|"tok(Token_ID)"| q154
    q154 --> q157
    q155 --> q153
    q155 --> q156
    q156 --> q31
    q157 --> q155
```

## J

```mermaid
flowchart TD
    q32(["J__Start (32)<br/>RuleStart"])
    q33(["J__Stop (33)<br/>RuleStop"])
    q158["J_j (158)<br/>Basic<br/>"]
    q159["J_Ref_ID (159)<br/>Basic<br/>"]
    q160["J__Basic_0 (160)<br/>Basic<br/>"]
    q161["J_SELF (161)<br/>Basic<br/>"]
    q162["J__Basic_1 (162)<br/>Basic<br/>"]
    q163{"J__Basic_2 (163)<br/>Basic<br/><br/>dec=10"}
    q164["J__BlockEnd (164)<br/>BlockEnd<br/>"]

    q32 --> q158
    q158 -->|"tok(Keyword_j)"| q163
    q159 -->|"tok(Token_ID)"| q160
    q160 --> q164
    q161 -->|"tok(Token_SELF)"| q162
    q162 --> q164
    q163 --> q159
    q163 --> q161
    q164 --> q33
```

## K

```mermaid
flowchart TD
    q34(["K__Start (34)<br/>RuleStart"])
    q35(["K__Stop (35)<br/>RuleStop"])
    q165["K_k (165)<br/>Basic<br/>"]
    q166["K_Ref1_ID (166)<br/>Basic<br/>"]
    q167["K_x (167)<br/>Basic<br/>"]
    q168["K__Basic_0 (168)<br/>Basic<br/>"]
    q169["K_Ref2_ID (169)<br/>Basic<br/>"]
    q170["K_y (170)<br/>Basic<br/>"]
    q171["K__Basic_1 (171)<br/>Basic<br/>"]
    q172{"K__Basic_2 (172)<br/>Basic<br/><br/>dec=11"}
    q173["K__BlockEnd (173)<br/>BlockEnd<br/>"]

    q34 --> q165
    q165 -->|"tok(Keyword_k)"| q172
    q166 -->|"tok(Token_ID)"| q167
    q167 -->|"tok(Keyword_x)"| q168
    q168 --> q173
    q169 -->|"tok(Token_ID)"| q170
    q170 -->|"tok(Keyword_y)"| q171
    q171 --> q173
    q172 --> q166
    q172 --> q169
    q173 --> q35
```

## L

```mermaid
flowchart TD
    q36(["L__Start (36)<br/>RuleStart"])
    q37(["L__Stop (37)<br/>RuleStop"])
    q174["L_l (174)<br/>Basic<br/>"]
    q175["L_OPTIONAL (175)<br/>Basic<br/>"]
    q176["L_AND (176)<br/>Basic<br/>"]
    q177["L__Basic_0 (177)<br/>Basic<br/>"]
    q178{"L__Basic_1 (178)<br/>Basic<br/><br/>dec=12"}
    q179["L_THEN (179)<br/>Basic<br/>"]
    q180["L_END (180)<br/>Basic<br/>"]
    q181["L__Basic_2 (181)<br/>Basic<br/>"]

    q36 --> q174
    q174 -->|"tok(Keyword_l)"| q178
    q175 -->|"tok(Token_OPTIONAL)"| q176
    q176 -->|"tok(Token_AND)"| q177
    q177 --> q179
    q178 --> q175
    q178 --> q177
    q179 -->|"tok(Token_THEN)"| q180
    q180 -->|"tok(Token_END)"| q181
    q181 --> q37
```

## M

```mermaid
flowchart TD
    q38(["M__Start (38)<br/>RuleStart"])
    q39(["M__Stop (39)<br/>RuleStop"])
    q182["M_m (182)<br/>Basic<br/>"]
    q183["M__Basic_0 (183)<br/>Basic<br/>"]
    q184["M__Basic_1 (184)<br/>Basic<br/>"]

    q38 --> q182
    q182 -->|"tok(Keyword_m)"| q183
    q183 -->|"tok(TokenGroup_SomeTokenGroup)"| q184
    q184 --> q39
```

## N

```mermaid
flowchart TD
    q40(["N__Start (40)<br/>RuleStart"])
    q41(["N__Stop (41)<br/>RuleStop"])
    q185["N_n (185)<br/>Basic<br/>"]
    q186["N__Basic_0 (186)<br/>Basic<br/>"]
    q187["N__Basic_1 (187)<br/>Basic<br/>"]

    q40 --> q185
    q185 -->|"tok(Keyword_n)"| q186
    q186 -->|"tok(TokenGroup_SomeTokenGroup)"| q187
    q187 --> q41
```

## O

```mermaid
flowchart TD
    q42(["O__Start (42)<br/>RuleStart"])
    q43(["O__Stop (43)<br/>RuleStop"])
    q188["O_o (188)<br/>Basic<br/>"]
    q189["O_Ref_ID (189)<br/>Basic<br/>"]
    q190["O__Basic (190)<br/>Basic<br/>"]

    q42 --> q188
    q188 -->|"tok(Keyword_o)"| q189
    q189 -->|"tok(Token_ID)"| q190
    q190 --> q43
```

## FQN

```mermaid
flowchart TD
    q44(["FQN__Start (44)<br/>RuleStart"])
    q45(["FQN__Stop (45)<br/>RuleStop"])
    q191["FQN_ID_0 (191)<br/>Basic<br/>"]
    q192["FQN_DOT (192)<br/>Basic<br/>"]
    q193["FQN_ID_1 (193)<br/>Basic<br/>"]
    q194["FQN__Basic (194)<br/>Basic<br/>"]
    q195{"FQN__LoopEntry (195)<br/>LoopEntry<br/><br/>dec=13"}
    q196["FQN__LoopEnd (196)<br/>LoopEnd<br/>"]
    q197["FQN__LoopBack (197)<br/>LoopBack<br/>"]

    q44 --> q191
    q191 -->|"tok(Token_ID)"| q195
    q192 -->|"tok(Token_DOT)"| q193
    q193 -->|"tok(Token_ID)"| q194
    q194 --> q197
    q195 --> q192
    q195 --> q196
    q196 --> q45
    q197 --> q195
```

