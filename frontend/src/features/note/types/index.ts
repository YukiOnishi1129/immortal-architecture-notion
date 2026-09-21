import type { NOTE_STATUS } from "../constants";

export type NoteStatus = (typeof NOTE_STATUS)[keyof typeof NOTE_STATUS];

export interface NoteSection {
  id: string;
  fieldId: string;
  fieldLabel: string;
  content: string;
  isRequired: boolean;
}

export interface NoteOwner {
  id: string;
  firstName: string;
  lastName: string;
  thumbnail: string | null;
}

export interface Note {
  id: string;
  title: string;
  templateId: string;
  templateName: string;
  ownerId: string;
  owner: NoteOwner;
  status: NoteStatus;
  sections: NoteSection[];
  createdAt: string;
  updatedAt: string;
}

export interface NoteFilters {
  status?: NoteStatus;
  templateId?: string;
  search?: string;
  q?: string;
  page?: number;
  ownerId?: string;
}
