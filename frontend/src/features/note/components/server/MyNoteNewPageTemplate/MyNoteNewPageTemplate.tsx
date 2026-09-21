import { NoteNewForm } from "@/features/note/components/client/NoteNewForm";

export function MyNoteNewPageTemplate() {
  return (
    <div className="container mx-auto py-6">
      <NoteNewForm backTo="/my-notes" />
    </div>
  );
}
