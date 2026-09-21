export const NOTE_STATUS = {
  DRAFT: "Draft",
  PUBLISH: "Publish",
} as const;

export type NoteStatusType = (typeof NOTE_STATUS)[keyof typeof NOTE_STATUS];

export const NOTE_STATUS_LABELS = {
  [NOTE_STATUS.DRAFT]: "下書き",
  [NOTE_STATUS.PUBLISH]: "公開",
} as const;
