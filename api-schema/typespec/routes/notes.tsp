import "@typespec/http";
import "@typespec/openapi3";
import "../models/note.tsp";
import "../models/common.tsp";

using TypeSpec.Http;
using MiniNotion.Models;

namespace MiniNotion.Routes;

@route("/api/notes")
@tag("Notes")
interface Notes {
  /** ノート一覧取得 */
  @get
  @summary("Get notes list")
  listNotes(
    /** タイトルキーワード検索 */
    @query q?: string,

    /** ステータスフィルター */
    @query status?: NoteStatus,

    /** テンプレートIDフィルター */
    @query templateId?: string,

    /** 所有者IDフィルター */
    @query ownerId?: string
  ): NoteResponse[] | UnauthorizedError;

  /** ノート詳細取得 */
  @get
  @route("/{noteId}")
  @summary("Get note by ID")
  getNoteById(
    @path noteId: string
  ): NoteResponse | NotFoundError | UnauthorizedError;

  /** ノート作成 */
  @post
  @summary("Create note")
  createNote(
    @body request: CreateNoteRequest
  ): NoteResponse | BadRequestError | UnauthorizedError;

  /** ノート更新 */
  @put
  @route("/{noteId}")
  @summary("Update note")
  updateNote(
    @path noteId: string,
    /** 所有者ID（権限チェック用） */
    @query ownerId: string,
    @body request: UpdateNoteRequest
  ): NoteResponse | NotFoundError | ForbiddenError | BadRequestError | UnauthorizedError;

  /** ノート公開 */
  @post
  @route("/{noteId}/publish")
  @summary("Publish note")
  publishNote(
    @path noteId: string,
    /** 所有者ID（公開権限チェック用） */
    @query ownerId: string
  ): NoteResponse | NotFoundError | ForbiddenError | BadRequestError | UnauthorizedError;

  /** ノート公開取り消し */
  @post
  @route("/{noteId}/unpublish")
  @summary("Unpublish note")
  unpublishNote(
    @path noteId: string,
    /** 所有者ID（公開権限チェック用） */
    @query ownerId: string
  ): NoteResponse | NotFoundError | ForbiddenError | BadRequestError | UnauthorizedError;

  /** ノート削除 */
  @delete
  @route("/{noteId}")
  @summary("Delete note")
  deleteNote(
    @path noteId: string,
    /** 所有者ID（権限チェック用） */
    @query ownerId: string
  ): SuccessResponse | NotFoundError | ForbiddenError | UnauthorizedError;
}
