use mysql_async::prelude::*;
use mysql_async::{FromRowError, OptsBuilder, Params, Row, Pool};
use mysql_common::row::ColumnIndex;

use mysql_common::value::Value;
use std::collections::HashMap;

pub struct XXX{}

///////////////// SHARED CODE ///////////
#[derive(Debug)]
pub struct SPool {
    pub pool: mysql_async::Pool,
    pub database: String,
}

#[derive(Default, Debug)]
pub struct TQuery {
    pub wheres: Vec<WhereClause>,
    pub wheres_ins: Vec<WhereInClause>,
    pub select_cols: Vec<&'static str>, // Just Selector
    pub delete_cols: Vec<&'static str>, // Deleter
    pub order_by:  Vec<&'static str>,
    pub updates: HashMap<&'static str, Value>,
    pub limit: u32,
    pub offset: u32,
}

#[derive(Debug, Clone)]
pub struct WhereClause {
    // pub condition: &'static str,
    pub condition: String,
    pub args: Value,
}

#[derive(Debug, Clone)]
pub struct WhereInClause {
    pub condition: String,
    pub args: Vec<Value>,
}

pub fn _get_where(wheres: Vec<WhereClause>) ->  (String, Vec<Value>) {
    let mut values = vec![];
    let  mut where_str = vec![];

    for w in wheres {
        where_str.push(w.condition);
        values.push(w.args)
    }
    let cql_where = where_str.join(" ");

    (cql_where, values)
}

impl TQuery {
    pub fn _to_sql_selector(&self) ->  (String, Vec<Value>)  {
        let cql_select = if self.select_cols.is_empty() {
            "*".to_string()
        } else {
            self.select_cols.join(", ")
        };

        let mut cql_query = format!("SELECT {} FROM `twitter`.`tweet`", cql_select);

        let (cql_where, where_values) = _get_where(self.wheres.clone());

        if where_values.len() > 0 {
            cql_query.push_str(&format!(" WHERE {}",&cql_where));
        }

        if self.order_by.len() > 0 {
            let cql_orders = self.order_by.join(", ");
            cql_query.push_str( &format!(" ORDER BY {}", &cql_orders));
        };

        if self.limit != 0  {
            cql_query.push_str(&format!(" LIMIT {} ", self.limit));
        };

        if self.offset != 0  {
            cql_query.push_str(&format!(" OFFSET {} ", self.offset));
        };

        (cql_query, where_values)
    }

}

pub async fn update_rows(query: &TQuery, session: &SPool) -> Result<(),MyError> {
    let mut conn = session.pool.get_conn().await?;

    if query.updates.is_empty() {
        return Err(MyError::EmptySql);
    }

    // Update columns building
    let mut all_vals = vec![];
    let mut col_updates = vec![];

    for (col,val) in query.updates.clone() {
        all_vals.push(val);
        col_updates.push(col);
    }
    let cql_update = col_updates.join(",");

    // Where columns building
    let  mut where_str = vec![];

    for w in query.wheres.clone() {
        where_str.push(w.condition);
        all_vals.push(w.args)
    }
    let cql_where = where_str.join(" ");

    // Build final query
    let mut cql_query = if query.wheres.is_empty() {
        format!("UPDATE .TableSchemeOut SET {}", cql_update)
    } else {
        format!("UPDATE .TableSchemeOut SET {} WHERE {}", cql_update, cql_where)
    };

    let p = Params::Positional(all_vals);

    println!("{} - {:?}", &cql_query, &p);

    let query_result = conn.exec_drop(cql_query,p,).await?;

    Ok(())
}

pub async fn delete_rows(query: &TQuery, session: &SPool) -> Result<(),MyError> {
    let mut conn = session.pool.get_conn().await?;
    let del_col = query.delete_cols.join(", ");

    let (cql_where, where_values) = _get_where(query.wheres.clone());

    let cql_query = format!("DELETE {} FROM `twitter`.`tweet` WHERE {}", del_col, cql_where);

    let p = Params::Positional(where_values);

    println!("{} - {:?}", &cql_query, &p);

    let query_result = conn.exec_drop(cql_query,p,).await?;

    Ok(())
}


#[derive(Debug)]
pub enum MyError { // MySQL Error
    NotFound,
    EmptySql,
    MySqlError(mysql_async::Error),
}

impl From<mysql_async::Error> for MyError{
    fn from(err: mysql_async::Error) -> Self {
        MyError::MySqlError(err)
    }
}

